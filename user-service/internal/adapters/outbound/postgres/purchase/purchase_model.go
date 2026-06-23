package postgrespurchase

import (
	"time"
	"user-service/internal/domain"
)

type PurchaseModel struct {
	ID          int64               `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      int64               `gorm:"column:user_id;not null"`
	TotalAmount float64             `gorm:"column:total_amount;type:numeric(12,2);not null"`
	Items       []PurchaseItemModel `gorm:"foreignKey:PurchaseID"`
	CreatedAt   time.Time           `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time           `gorm:"column:updated_at;autoUpdateTime"`
}

func (PurchaseModel) TableName() string {
	return "purchases"
}

type PurchaseItemModel struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement"`
	PurchaseID      int64     `gorm:"column:purchase_id;not null"`
	ClothesID       int64     `gorm:"column:clothes_id;not null"`
	ProductName     string    `gorm:"column:product_name;size:180;not null"`
	ProductImageURL *string   `gorm:"column:product_image_url"`
	Quantity        int       `gorm:"column:quantity;not null"`
	UnitPrice       float64   `gorm:"column:unit_price;type:numeric(12,2);not null"`
	Size            string    `gorm:"column:size;size:20;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (PurchaseItemModel) TableName() string {
	return "purchase_items"
}

func toDomain(model PurchaseModel) domain.Purchase {
	items := make([]domain.PurchaseItem, 0, len(model.Items))
	for _, item := range model.Items {
		items = append(items, domain.PurchaseItem{
			ID:              item.ID,
			PurchaseID:      item.PurchaseID,
			ClothesID:       item.ClothesID,
			ProductName:     item.ProductName,
			ProductImageURL: item.ProductImageURL,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			Size:            item.Size,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
		})
	}

	return domain.Purchase{
		ID:          model.ID,
		UserID:      model.UserID,
		TotalAmount: model.TotalAmount,
		Items:       items,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
