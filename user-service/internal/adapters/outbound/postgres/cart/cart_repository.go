package postgrescart

import (
	"context"
	"errors"
	"user-service/internal/application/ports"
	"user-service/internal/domain"

	"gorm.io/gorm"
)

var _ ports.CartRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindActiveByUserID(ctx context.Context, userID int64) (*domain.Cart, error) {
	return findActiveByUserID(ctx, r.db, userID)
}

func (r *Repository) AddItem(ctx context.Context, userID int64, item domain.CartItem) (*domain.Cart, error) {
	var cart *domain.Cart

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cartModel, err := findOrCreateActiveCart(ctx, tx, userID)
		if err != nil {
			return err
		}

		var itemModel CartItemModel
		err = tx.WithContext(ctx).
			Where("cart_id = ? AND clothes_id = ? AND size = ?", cartModel.ID, item.ClothesID, item.Size).
			Take(&itemModel).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			itemModel = CartItemModel{
				CartID:    cartModel.ID,
				ClothesID: item.ClothesID,
				Quantity:  item.Quantity,
				Size:      item.Size,
			}
			if err := tx.WithContext(ctx).Create(&itemModel).Error; err != nil {
				return err
			}
		} else {
			itemModel.Quantity += item.Quantity
			if err := tx.WithContext(ctx).Save(&itemModel).Error; err != nil {
				return err
			}
		}

		cart, err = findActiveByUserID(ctx, tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *Repository) UpdateItemQuantity(ctx context.Context, userID int64, itemID int64, quantity int) (*domain.Cart, error) {
	var cart *domain.Cart

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cartModel, err := findActiveCartModel(ctx, tx, userID)
		if err != nil {
			return err
		}

		result := tx.WithContext(ctx).
			Model(&CartItemModel{}).
			Where("id = ? AND cart_id = ?", itemID, cartModel.ID).
			Update("quantity", quantity)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrCartItemNotFound
		}

		cart, err = findActiveByUserID(ctx, tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *Repository) DeleteItem(ctx context.Context, userID int64, itemID int64) (*domain.Cart, error) {
	var cart *domain.Cart

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cartModel, err := findActiveCartModel(ctx, tx, userID)
		if err != nil {
			return err
		}

		result := tx.WithContext(ctx).
			Where("id = ? AND cart_id = ?", itemID, cartModel.ID).
			Delete(&CartItemModel{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrCartItemNotFound
		}

		cart, err = findActiveByUserID(ctx, tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *Repository) Clear(ctx context.Context, userID int64) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, domain.CartStatusActive).
		Delete(&CartModel{})
	return result.Error
}

func findOrCreateActiveCart(ctx context.Context, db *gorm.DB, userID int64) (*CartModel, error) {
	cart, err := findActiveCartModel(ctx, db, userID)
	if err == nil {
		return cart, nil
	}
	if !errors.Is(err, domain.ErrCartNotFound) {
		return nil, err
	}

	model := &CartModel{
		UserID: userID,
		Status: string(domain.CartStatusActive),
	}
	if err := db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func findActiveByUserID(ctx context.Context, db *gorm.DB, userID int64) (*domain.Cart, error) {
	model, err := findActiveCartModel(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	if err := db.WithContext(ctx).
		Where("cart_id = ?", model.ID).
		Order("id ASC").
		Find(&model.Items).Error; err != nil {
		return nil, err
	}

	cart := toDomain(*model)
	return &cart, nil
}

func findActiveCartModel(ctx context.Context, db *gorm.DB, userID int64) (*CartModel, error) {
	var model CartModel
	if err := db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, domain.CartStatusActive).
		Take(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCartNotFound
		}

		return nil, err
	}

	return &model, nil
}
