package ports

import (
	"context"

	"clothing-service/internal/domain"
)

type ClothesImageRepository interface {
	Create(ctx context.Context, image *domain.ClothesImage) error
	ListByClothesID(ctx context.Context, clothesID int64) ([]domain.ClothesImage, error)
	FindByID(ctx context.Context, imageID int64) (*domain.ClothesImage, error)
	Delete(ctx context.Context, imageID int64) error
	DeleteByClothesID(ctx context.Context, clothesID int64) error
}
