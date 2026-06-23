package postgrescolor

import "clothing-service/internal/domain"

func ToColorModel(color *domain.Color) *ColorModel {
	if color == nil {
		return nil
	}

	return &ColorModel{
		ID:        color.ID,
		Name:      color.Name,
		HexCode:   color.HexCode,
		CreatedAt: color.CreatedAt,
		UpdatedAt: color.UpdatedAt,
	}
}

func ToColorDomain(model ColorModel) domain.Color {
	return domain.Color{
		ID:        model.ID,
		Name:      model.Name,
		HexCode:   model.HexCode,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
