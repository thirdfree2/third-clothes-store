package postgrescolor

import (
	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"
	"context"

	"gorm.io/gorm"
)

var _ ports.ColorRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllColor(ctx context.Context) ([]domain.Color, error) {
	var models []ColorModel

	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	colors := make([]domain.Color, 0, len(models))

	for _, model := range models {
		colors = append(colors, domain.Color{
			ID:        model.ID,
			Name:      model.Name,
			HexCode:   model.HexCode,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		})
	}

	return colors, nil
}
