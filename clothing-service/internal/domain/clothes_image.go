package domain

import "time"

type ClothesImage struct {
	ID               int64
	ClothesID        int64
	Bucket           string
	ObjectKey        string
	OriginalFilename *string
	ContentType      string
	SizeBytes        int64
	CreatedAt        time.Time
}
