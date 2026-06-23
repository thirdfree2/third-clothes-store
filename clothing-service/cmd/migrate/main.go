package main

import (
	postgresdb "clothing-service/internal/adapters/outbound/postgres/db"
	"clothing-service/internal/config"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

const migrationsDir = "migrations"

type migrationFile struct {
	Version string
	Path    string
}

func main() {
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	cfg := config.Load()
	db, err := postgresdb.NewDB(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	switch command {
	case "up":
		err = migrateUp(db)
	case "down":
		err = migrateDown(db)
	case "status":
		err = showStatus(db)
	default:
		err = fmt.Errorf("unknown command %q, use up, down, or status", command)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func migrateUp(db *gorm.DB) error {
	if err := ensureSchemaMigrations(db); err != nil {
		return err
	}

	files, err := migrationFiles(".up.sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		applied, err := isApplied(db, file.Version)
		if err != nil {
			return err
		}
		if applied {
			log.Printf("skip %s", filepath.Base(file.Path))
			continue
		}

		if err := applyUp(db, file); err != nil {
			return err
		}
		log.Printf("applied %s", filepath.Base(file.Path))
	}

	return nil
}

func migrateDown(db *gorm.DB) error {
	if err := ensureSchemaMigrations(db); err != nil {
		return err
	}

	version, err := latestAppliedVersion(db)
	if err != nil {
		return err
	}
	if version == "" {
		log.Println("no migration to rollback")
		return nil
	}

	path := filepath.Join(migrationsDir, version+".down.sql")
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("down migration not found: %s", path)
		}
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(string(content)).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", version).Error; err != nil {
			return err
		}

		log.Printf("rolled back %s", filepath.Base(path))
		return nil
	})
}

func showStatus(db *gorm.DB) error {
	if err := ensureSchemaMigrations(db); err != nil {
		return err
	}

	files, err := migrationFiles(".up.sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		applied, err := isApplied(db, file.Version)
		if err != nil {
			return err
		}

		status := "pending"
		if applied {
			status = "applied"
		}
		log.Printf("%s %s", file.Version, status)
	}

	return nil
}

func applyUp(db *gorm.DB, file migrationFile) error {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(string(content)).Error; err != nil {
			return err
		}
		return tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", file.Version).Error
	})
}

func ensureSchemaMigrations(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version text PRIMARY KEY,
			applied_at timestamptz NOT NULL DEFAULT now()
		)
	`).Error
}

func isApplied(db *gorm.DB, version string) (bool, error) {
	var count int64
	err := db.Raw("SELECT count(*) FROM schema_migrations WHERE version = ?", version).Scan(&count).Error
	return count > 0, err
}

func latestAppliedVersion(db *gorm.DB) (string, error) {
	var version string
	err := db.Raw("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version).Error
	return version, err
}

func migrationFiles(suffix string) ([]migrationFile, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	files := make([]migrationFile, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), suffix) {
			continue
		}

		version := strings.TrimSuffix(entry.Name(), suffix)
		files = append(files, migrationFile{
			Version: version,
			Path:    filepath.Join(migrationsDir, entry.Name()),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	return files, nil
}
