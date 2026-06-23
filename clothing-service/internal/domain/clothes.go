package domain

import "time"

type Clothes struct {
	ID         int64
	Name       string
	Price      float64
	ColorID    *int64
	Color      *Color
	Categories []Category
	Images     []ClothesImage
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
