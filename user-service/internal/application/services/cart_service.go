package services

import (
	"context"
	"strings"
	"user-service/internal/application/ports"
	"user-service/internal/domain"
)

type CartService struct {
	cartRepo ports.CartRepository
}

func NewCartService(cartRepo ports.CartRepository) *CartService {
	return &CartService{cartRepo: cartRepo}
}

func (s *CartService) GetCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	cart, err := s.cartRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		if err == domain.ErrCartNotFound {
			return &domain.Cart{
				UserID: userID,
				Status: domain.CartStatusActive,
				Items:  []domain.CartItem{},
			}, nil
		}

		return nil, err
	}

	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, userID int64, clothesID int64, quantity int, size string) (*domain.Cart, error) {
	if userID <= 0 || clothesID <= 0 || quantity <= 0 {
		return nil, domain.ErrInvalidCartItem
	}

	return s.cartRepo.AddItem(ctx, userID, domain.CartItem{
		ClothesID: clothesID,
		Quantity:  quantity,
		Size:      strings.TrimSpace(size),
	})
}

func (s *CartService) UpdateItemQuantity(ctx context.Context, userID int64, itemID int64, quantity int) (*domain.Cart, error) {
	if userID <= 0 || itemID <= 0 || quantity <= 0 {
		return nil, domain.ErrInvalidCartItem
	}

	return s.cartRepo.UpdateItemQuantity(ctx, userID, itemID, quantity)
}

func (s *CartService) DeleteItem(ctx context.Context, userID int64, itemID int64) (*domain.Cart, error) {
	if userID <= 0 || itemID <= 0 {
		return nil, domain.ErrInvalidCartItem
	}

	return s.cartRepo.DeleteItem(ctx, userID, itemID)
}

func (s *CartService) ClearCart(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return domain.ErrInvalidCartItem
	}

	return s.cartRepo.Clear(ctx, userID)
}
