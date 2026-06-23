package postgresclothes

import (
	postgrescategory "clothing-service/internal/adapters/outbound/postgres/category"
	postgrescolor "clothing-service/internal/adapters/outbound/postgres/color"
	"clothing-service/internal/domain"
)

func ToClothesModel(clothes *domain.Clothes) ClothesModel {
	if clothes == nil {
		return ClothesModel{}
	}

	model := ClothesModel{
		ID:        clothes.ID,
		Name:      clothes.Name,
		Price:     clothes.Price,
		ColorID:   clothes.ColorID,
		CreatedAt: clothes.CreatedAt,
		UpdatedAt: clothes.UpdatedAt,
	}

	if clothes.Color != nil {
		model.Color = postgrescolor.ToColorModel(clothes.Color)
	}

	if len(clothes.Categories) > 0 {
		model.Categories = make([]postgrescategory.CategoryModel, 0, len(clothes.Categories))
		for _, category := range clothes.Categories {
			model.Categories = append(model.Categories, postgrescategory.ToCategoryModel(&category))
		}
	}

	return model
}

func ToClothesDomain(model ClothesModel) domain.Clothes {
	clothes := domain.Clothes{
		ID:         model.ID,
		Name:       model.Name,
		Price:      model.Price,
		ColorID:    model.ColorID,
		Categories: make([]domain.Category, 0, len(model.Categories)),
		Images:     make([]domain.ClothesImage, 0, len(model.Images)),
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}

	if model.Color != nil {
		clothes.Color = &domain.Color{
			ID:      model.Color.ID,
			Name:    model.Color.Name,
			HexCode: model.Color.HexCode,
		}
	}

	for _, category := range model.Categories {
		clothes.Categories = append(clothes.Categories, domain.Category{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	for _, image := range model.Images {
		clothes.Images = append(clothes.Images, domain.ClothesImage{
			ID:               image.ID,
			ClothesID:        image.ClothesID,
			Bucket:           image.Bucket,
			ObjectKey:        image.ObjectKey,
			OriginalFilename: image.OriginalFilename,
			ContentType:      image.ContentType,
			SizeBytes:        image.SizeBytes,
			CreatedAt:        image.CreatedAt,
		})
	}

	return clothes
}
