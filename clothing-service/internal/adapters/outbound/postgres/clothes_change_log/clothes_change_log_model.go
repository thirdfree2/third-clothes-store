package postgresclotheschangelog

import (
	"time"

	"gorm.io/datatypes"
)

type clothesChangeLogModel struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement"`

	ClothesID *int64 `gorm:"column:clothes_id"`

	Action string `gorm:"column:action"`

	BeforeData datatypes.JSON `gorm:"column:before_data;type:jsonb"`
	AfterData  datatypes.JSON `gorm:"column:after_data;type:jsonb"`

	ChangedBy *string `gorm:"column:changed_by"`
	Note      *string `gorm:"column:note"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (clothesChangeLogModel) TableName() string {
	return "clothes_change_logs"
}
