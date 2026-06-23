package services

import "clothing-service/internal/domain"

func clothesToLogData(c *domain.Clothes) map[string]any {
	if c == nil {
		return nil
	}

	categories := make([]map[string]any, 0, len(c.Categories))
	categoryIDs := make([]int64, 0, len(c.Categories))

	for _, category := range c.Categories {
		categoryIDs = append(categoryIDs, category.ID)

		categories = append(categories, map[string]any{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	data := map[string]any{
		"id":           c.ID,
		"name":         c.Name,
		"price":        c.Price,
		"color_id":     c.ColorID,
		"category_ids": categoryIDs,
		"categories":   categories,
		"created_at":   c.CreatedAt,
		"updated_at":   c.UpdatedAt,
	}

	if c.Color != nil {
		data["color"] = map[string]any{
			"id":   c.Color.ID,
			"name": c.Color.Name,
		}
	}

	return data
}

func clothesImagesToLogData(images []domain.ClothesImage) []map[string]any {
	result := make([]map[string]any, 0, len(images))

	for _, image := range images {
		result = append(result, map[string]any{
			"id":                image.ID,
			"clothes_id":        image.ClothesID,
			"bucket":            image.Bucket,
			"object_key":        image.ObjectKey,
			"original_filename": image.OriginalFilename,
			"content_type":      image.ContentType,
			"size_bytes":        image.SizeBytes,
			"created_at":        image.CreatedAt,
		})
	}

	return result
}
