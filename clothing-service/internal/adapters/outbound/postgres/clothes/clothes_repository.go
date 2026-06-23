package postgresclothes

import (
	"context"
	"errors"

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
) ([]domain.Clothes, int64, error) {
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&ClothesModel{}).
		Where("deleted_at IS NULL")

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []ClothesModel

	if err := baseQuery.
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
