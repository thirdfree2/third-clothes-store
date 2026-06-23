package domain

import "time"

type Purchase struct {
	ID          int64
	UserID      int64
	TotalAmount float64
	Items       []PurchaseItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PurchaseItem struct {
	ID              int64
	PurchaseID      int64
	ClothesID       int64
	ProductName     string
	ProductImageURL *string
	Quantity        int
	UnitPrice       float64
	Size            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
