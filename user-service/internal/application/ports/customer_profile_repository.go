package ports

import (
	"context"
	"user-service/internal/domain"
)

type CustomerProfileRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*domain.CustomerProfile, error)
	Upsert(ctx context.Context, profile *domain.CustomerProfile) error
}

type CustomerAddressRepository interface {
	FindDefaultByUserID(ctx context.Context, userID int64) (*domain.CustomerAddress, error)
	UpsertDefault(ctx context.Context, address *domain.CustomerAddress) error
}
