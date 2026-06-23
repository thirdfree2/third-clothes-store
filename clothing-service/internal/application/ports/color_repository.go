package ports

import (
	"clothing-service/internal/domain"
	"context"
)

type ColorRepository interface {
	GetAllColor(ctx context.Context) ([]domain.Color, error)
}
