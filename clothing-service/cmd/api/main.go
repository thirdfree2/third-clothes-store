package main

import (
	httprouter "clothing-service/internal/adapters/inbound/http/router"
	minioadapter "clothing-service/internal/adapters/outbound/minio"
	postgrescategory "clothing-service/internal/adapters/outbound/postgres/category"
	postgresclothes "clothing-service/internal/adapters/outbound/postgres/clothes"
	postgresclotheschangelog "clothing-service/internal/adapters/outbound/postgres/clothes_change_log"
	postgresclothesimage "clothing-service/internal/adapters/outbound/postgres/clothes_image"
	postgrescolor "clothing-service/internal/adapters/outbound/postgres/color"
	postgresdb "clothing-service/internal/adapters/outbound/postgres/db"
	"clothing-service/internal/application/services"
	"clothing-service/internal/config"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

// @title Go Gin Hexagonal API
// @version 1.0
// @description REST API example using Gin, Hexagonal Architecture, PostgreSQL, and GORM.
// @host localhost:9080
// @BasePath /
func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := postgresdb.NewDB(cfg.DatabaseDSN())
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	objectStorage, err := minioadapter.NewObjectStorage(
		cfg.MinIOEndpoint,
		cfg.MinIOAccessKey,
		cfg.MinIOSecretKey,
		cfg.MinIOUseSSL,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := objectStorage.EnsureBucket(
		context.Background(),
		cfg.MinIOBucket,
	); err != nil {
		log.Fatal(err)
	}

	clothesRepo := postgresclothes.NewRepository(db)
	clothesChangeLogRepo := postgresclotheschangelog.NewRepository(db)
	clothesImageRepo := postgresclothesimage.NewRepository(db)
	colorRepo := postgrescolor.NewRepository(db)
	categoryRepo := postgrescategory.NewRepository(db)

	clothesImageService := services.NewClothesImageService(clothesRepo, clothesImageRepo, objectStorage, clothesChangeLogRepo, cfg.MinIOBucket)
	clothesService := services.NewClothesService(clothesRepo, clothesImageRepo, clothesChangeLogRepo)
	colorService := services.NewColorService(colorRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	router := httprouter.NewRouter(clothesService, clothesImageService, colorService, categoryService, cfg.MinIOPublicURL, cfg.JWTSecret,
		cfg.JWTIssuer)

	log.Printf("server listening on %s", cfg.HTTPAddr())
	if err := router.Run(cfg.HTTPAddr()); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
