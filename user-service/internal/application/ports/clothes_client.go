package ports

import (
	"context"
	"user-service/internal/domain"
)

type ClothesClient interface {
	GetByID(ctx context.Context, id int64) (*domain.Clothes, error)
}
