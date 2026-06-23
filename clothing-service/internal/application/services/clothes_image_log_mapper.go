package services

import "clothing-service/internal/domain"

func clothesImageToLogData(image *domain.ClothesImage) map[string]any {
	if image == nil {
		return nil
	}

	return map[string]any{
		"id":                image.ID,
		"clothes_id":        image.ClothesID,
		"bucket":            image.Bucket,
		"object_key":        image.ObjectKey,
		"original_filename": image.OriginalFilename,
		"content_type":      image.ContentType,
		"size_bytes":        image.SizeBytes,
		"created_at":        image.CreatedAt,
	}
}
