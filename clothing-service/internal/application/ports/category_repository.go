package ports

import (
	"clothing-service/internal/domain"
	"context"
)

type CategoryRepository interface {
	GetAllCategory(ctx context.Context) ([]domain.Category, error)
}
