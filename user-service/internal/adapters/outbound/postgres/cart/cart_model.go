package postgrescart

import (
	"time"
	"user-service/internal/domain"
)

type CartModel struct {
	ID        int64           `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64           `gorm:"column:user_id;not null"`
	Status    string          `gorm:"column:status;size:20;not null"`
	Items     []CartItemModel `gorm:"foreignKey:CartID"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime"`
}

func (CartModel) TableName() string {
	return "carts"
}

type CartItemModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CartID    int64     `gorm:"column:cart_id;not null"`
	ClothesID int64     `gorm:"column:clothes_id;not null"`
	Quantity  int       `gorm:"column:quantity;not null"`
	Size      string    `gorm:"column:size;size:20;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (CartItemModel) TableName() string {
	return "cart_items"
}

func toDomain(model CartModel) domain.Cart {
	items := make([]domain.CartItem, 0, len(model.Items))
	for _, item := range model.Items {
		items = append(items, domain.CartItem{
			ID:        item.ID,
			CartID:    item.CartID,
			ClothesID: item.ClothesID,
			Quantity:  item.Quantity,
			Size:      item.Size,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return domain.Cart{
		ID:        model.ID,
		UserID:    model.UserID,
		Status:    domain.CartStatus(model.Status),
		Items:     items,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
