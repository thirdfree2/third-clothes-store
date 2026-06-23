package ports

import (
	"context"
	"user-service/internal/domain"
)

type CustomerProfileRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*domain.CustomerProfile, error)
	Upsert(ctx context.Context, profile *domain.CustomerProfile) error
}
