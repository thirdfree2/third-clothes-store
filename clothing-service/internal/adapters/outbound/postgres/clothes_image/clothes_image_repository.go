package postgresclothesimage

import (
	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

var _ ports.ClothesImageRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, image *domain.ClothesImage) error {
	model := ClothesImageModel{
		ClothesID:        image.ClothesID,
		Bucket:           image.Bucket,
		ObjectKey:        image.ObjectKey,
		OriginalFilename: image.OriginalFilename,
		ContentType:      image.ContentType,
		SizeBytes:        image.SizeBytes,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	image.ID = model.ID
	image.CreatedAt = model.CreatedAt

	return nil
}

func (r *Repository) FindByID(ctx context.Context, imageID int64) (*domain.ClothesImage, error) {
	var model ClothesImageModel

	if err := r.db.WithContext(ctx).First(&model, imageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrClothesImageNotFound
		}
		return nil, err
	}

	image := toDomain(model)
	return &image, nil
}

func (r *Repository) ListByClothesID(ctx context.Context, clothesID int64) ([]domain.ClothesImage, error) {
	var models []ClothesImageModel

	if err := r.db.WithContext(ctx).
		Where("clothes_id = ?", clothesID).
		Order("id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	images := make([]domain.ClothesImage, 0, len(models))
	for _, model := range models {
		images = append(images, toDomain(model))
	}

	return images, nil
}

func (r *Repository) Delete(ctx context.Context, imageID int64) error {
	var model ClothesImageModel

	if err := r.db.WithContext(ctx).First(&model, imageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrClothesImageNotFound
		}
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&model).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteByClothesID(ctx context.Context, clothesID int64) error {
	return r.db.WithContext(ctx).
		Exec(`
			UPDATE clothes_images
			SET deleted_at = NOW()
			WHERE clothes_id = ?
			  AND deleted_at IS NULL
		`, clothesID).Error
}
