package postgresclothesimage

import (
	"clothing-service/internal/domain"
	"time"

	"gorm.io/gorm"
)

type ClothesImageModel struct {
	ID               int64          `gorm:"column:id;primaryKey;autoIncrement"`
	ClothesID        int64          `gorm:"column:clothes_id"`
	Bucket           string         `gorm:"column:bucket"`
	ObjectKey        string         `gorm:"column:object_key"`
	OriginalFilename *string        `gorm:"column:original_filename"`
	ContentType      string         `gorm:"column:content_type"`
	SizeBytes        int64          `gorm:"column:size_bytes"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ClothesImageModel) TableName() string {
	return "clothes_images"
}

func toDomain(model ClothesImageModel) domain.ClothesImage {
	return domain.ClothesImage{
		ID:               model.ID,
		ClothesID:        model.ClothesID,
		Bucket:           model.Bucket,
		ObjectKey:        model.ObjectKey,
		OriginalFilename: model.OriginalFilename,
		ContentType:      model.ContentType,
		SizeBytes:        model.SizeBytes,
		CreatedAt:        model.CreatedAt,
	}
}
