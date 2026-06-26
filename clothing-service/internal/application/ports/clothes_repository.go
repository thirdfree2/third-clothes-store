package ports

import (
	"clothing-service/internal/domain"
	"context"
	"time"
)

type ClothesRepository interface {
	Create(ctx context.Context, clothes *domain.Clothes, categoryIDs []int64) error
	Update(ctx context.Context, clothes *domain.Clothes, categoryIDs []int64) error
	FindByID(ctx context.Context, id int64) (*domain.Clothes, error)
	List(ctx context.Context, pagination Pagination, filter ClothesListFilter) ([]domain.Clothes, int64, error)
	Delete(ctx context.Context, id int64) error
}

type ClothesListFilter struct {
	CategoryID   *int64
	ColorID      *int64
	Name         *string
	Price        *float64
	CategoryName *string
	CreatedDate  *time.Time
}
