package domain

import "time"

type CartStatus string

const (
	CartStatusActive     CartStatus = "ACTIVE"
	CartStatusCheckedOut CartStatus = "CHECKED_OUT"
)

type Cart struct {
	ID        int64
	UserID    int64
	Status    CartStatus
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ID        int64
	CartID    int64
	ClothesID int64
	Clothes   *Clothes
	Quantity  int
	Size      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
