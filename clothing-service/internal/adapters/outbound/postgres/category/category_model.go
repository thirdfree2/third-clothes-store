package postgrescategory

import "time"

type CategoryModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;size:120;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
