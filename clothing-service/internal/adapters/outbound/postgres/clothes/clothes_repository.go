package postgresclothes

import (
	"context"
	"errors"
	"strings"
	"time"

	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var _ ports.ClothesRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, clothes *domain.Clothes, categoryIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := ToClothesModel(clothes)

		// กัน GORM พยายาม create association เอง
		model.Color = nil
		model.Categories = nil

		if err := tx.Create(&model).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return domain.ErrClothesNameAlreadyExists
			}

			return err
		}

		if err := replaceCategories(tx, model.ID, categoryIDs); err != nil {
			return err
		}

		createdModel, err := findByIDTx(tx, model.ID)
		if err != nil {
			return err
		}

		*clothes = ToClothesDomain(*createdModel)
		return nil
	})
}

func (r *Repository) Update(ctx context.Context, clothes *domain.Clothes, categoryIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model ClothesModel

		if err := tx.First(&model, clothes.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrClothesNotFound
			}
			return err
		}

		model.Name = clothes.Name
		model.Price = clothes.Price
		model.ColorID = clothes.ColorID

		if err := tx.Save(&model).Error; err != nil {
			return err
		}

		if err := replaceCategories(tx, model.ID, categoryIDs); err != nil {
			return err
		}

		updatedModel, err := findByIDTx(tx, model.ID)
		if err != nil {
			return err
		}

		*clothes = ToClothesDomain(*updatedModel)
		return nil
	})
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*domain.Clothes, error) {
	model, err := findByIDTx(r.db.WithContext(ctx), id)
	if err != nil {
		return nil, err
	}

	clothes := ToClothesDomain(*model)
	return &clothes, nil
}

func (r *Repository) List(
	ctx context.Context,
	pagination ports.Pagination,
	filter ports.ClothesListFilter,
) ([]domain.Clothes, int64, error) {
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&ClothesModel{}).
		Where("deleted_at IS NULL")
	needsDistinct := false

	if filter.CategoryID != nil {
		baseQuery = baseQuery.
			Joins("JOIN clothes_categories ON clothes_categories.clothes_id = clothes.id").
			Where("clothes_categories.category_id = ?", *filter.CategoryID)
		needsDistinct = true
	}

	if filter.ColorID != nil {
		baseQuery = baseQuery.Where("clothes.color_id = ?", *filter.ColorID)
	}

	if filter.Name != nil && strings.TrimSpace(*filter.Name) != "" {
		baseQuery = baseQuery.Where("LOWER(clothes.name) LIKE ?", "%"+strings.ToLower(strings.TrimSpace(*filter.Name))+"%")
	}

	if filter.Price != nil {
		baseQuery = baseQuery.Where("clothes.price = ?", *filter.Price)
	}

	if filter.CategoryName != nil && strings.TrimSpace(*filter.CategoryName) != "" {
		baseQuery = baseQuery.
			Joins("JOIN clothes_categories AS search_clothes_categories ON search_clothes_categories.clothes_id = clothes.id").
			Joins("JOIN categories AS search_categories ON search_categories.id = search_clothes_categories.category_id").
			Where("LOWER(search_categories.name) LIKE ?", "%"+strings.ToLower(strings.TrimSpace(*filter.CategoryName))+"%")
		needsDistinct = true
	}

	if filter.CreatedDate != nil {
		start := beginningOfDay(*filter.CreatedDate)
		baseQuery = baseQuery.Where("clothes.created_at >= ? AND clothes.created_at < ?", start, start.AddDate(0, 0, 1))
	}

	countQuery := baseQuery
	if needsDistinct {
		countQuery = countQuery.Distinct("clothes.id")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []ClothesModel
	listQuery := baseQuery
	if needsDistinct {
		listQuery = listQuery.Distinct("clothes.*")
	}

	if err := listQuery.
		Preload("Color").
		Preload("Categories").
		Preload("Images").
		Order("id ASC").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	clothesList := make([]domain.Clothes, 0, len(models))
	for _, model := range models {
		clothesList = append(clothesList, ToClothesDomain(model))
	}

	return clothesList, total, nil
}

func beginningOfDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Model(&ClothesModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrClothesNotFound
	}

	return nil
}

func findByIDTx(tx *gorm.DB, id int64) (*ClothesModel, error) {
	var model ClothesModel

	if err := tx.
		Preload("Color").
		Preload("Categories").
		Preload("Images").
		First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrClothesNotFound
		}
		return nil, err
	}

	return &model, nil
}

func replaceCategories(tx *gorm.DB, clothesID int64, categoryIDs []int64) error {
	if err := tx.
		Where("clothes_id = ?", clothesID).
		Delete(&ClothesCategoryModel{}).Error; err != nil {
		return err
	}

	categoryIDs = uniqueInt64s(categoryIDs)
	if len(categoryIDs) == 0 {
		return nil
	}

	links := make([]ClothesCategoryModel, 0, len(categoryIDs))
	for _, categoryID := range categoryIDs {
		links = append(links, ClothesCategoryModel{
			ClothesID:  clothesID,
			CategoryID: categoryID,
		})
	}

	return tx.Create(&links).Error
}

func uniqueInt64s(values []int64) []int64 {
	seen := make(map[int64]struct{}, len(values))
	result := make([]int64, 0, len(values))

	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}
