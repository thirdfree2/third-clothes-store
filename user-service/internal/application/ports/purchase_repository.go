package ports

import (
	"context"
	"user-service/internal/domain"
)

type PurchaseRepository interface {
	Create(ctx context.Context, purchase domain.Purchase) (*domain.Purchase, error)
	FindByUserID(ctx context.Context, userID int64) ([]domain.Purchase, error)
}
