package ports

import (
	"clothing-service/internal/domain"
	"context"
)

type ClothesChangeLogRepository interface {
	Create(ctx context.Context, log *domain.ClothesChangeLog) error
}
