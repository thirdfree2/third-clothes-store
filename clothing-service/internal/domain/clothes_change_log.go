package domain

import "time"

type ClothesChangeAction string

const (
	ClothesChangeActionCreate ClothesChangeAction = "CREATE"
	ClothesChangeActionUpdate ClothesChangeAction = "UPDATE"
	ClothesChangeActionDelete ClothesChangeAction = "DELETE"

	ClothesChangeActionImageUpload     ClothesChangeAction = "IMAGE_UPLOAD"
	ClothesChangeActionImageDelete     ClothesChangeAction = "IMAGE_DELETE"
	ClothesChangeActionImageReplace    ClothesChangeAction = "IMAGE_REPLACE"
	ClothesChangeActionImageSetPrimary ClothesChangeAction = "IMAGE_SET_PRIMARY"
)

type ClothesChangeLog struct {
	ID        int64
	ClothesID *int64

	Action     ClothesChangeAction
	BeforeData map[string]any
	AfterData  map[string]any

	ChangedBy *string
	Note      *string

	CreatedAt time.Time
}
