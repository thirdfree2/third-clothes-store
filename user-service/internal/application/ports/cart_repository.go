package ports

import (
	"context"
	"user-service/internal/domain"
)

type CartRepository interface {
	FindActiveByUserID(ctx context.Context, userID int64) (*domain.Cart, error)
	AddItem(ctx context.Context, userID int64, item domain.CartItem) (*domain.Cart, error)
	UpdateItemQuantity(ctx context.Context, userID int64, itemID int64, quantity int) (*domain.Cart, error)
	DeleteItem(ctx context.Context, userID int64, itemID int64) (*domain.Cart, error)
	Clear(ctx context.Context, userID int64) error
}
