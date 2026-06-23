package postgresclothes

import (
	postgrescategory "clothing-service/internal/adapters/outbound/postgres/category"
	postgresclothesimage "clothing-service/internal/adapters/outbound/postgres/clothes_image"
	postgrescolor "clothing-service/internal/adapters/outbound/postgres/color"
	"time"

	"gorm.io/gorm"
)

type ClothesModel struct {
	ID      int64   `gorm:"column:id;primaryKey"`
	Name    string  `gorm:"column:name;size:180;not null"`
	Price   float64 `gorm:"column:price;type:numeric(12,2);not null"`
	ColorID *int64  `gorm:"column:color_id"`

	Color *postgrescolor.ColorModel `gorm:"foreignKey:ColorID;references:ID"`

	Categories []postgrescategory.CategoryModel `gorm:"many2many:clothes_categories;joinForeignKey:ClothesID;joinReferences:CategoryID"`

	Images []postgresclothesimage.ClothesImageModel `gorm:"foreignKey:ClothesID;references:ID"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ClothesModel) TableName() string {
	return "clothes"
}

type ClothesCategoryModel struct {
	ClothesID  int64 `gorm:"column:clothes_id;primaryKey"`
	CategoryID int64 `gorm:"column:category_id;primaryKey"`
}

func (ClothesCategoryModel) TableName() string {
	return "clothes_categories"
}
