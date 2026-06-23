package postgrescolor

import "time"

type ColorModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;size:120;not null"`
	HexCode   string    `gorm:"column:hex_code;size:20"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ColorModel) TableName() string {
	return "colors"
}
