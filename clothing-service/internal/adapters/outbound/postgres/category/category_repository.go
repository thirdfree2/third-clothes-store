package postgrescategory

import (
	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"
	"context"

	"gorm.io/gorm"
)

var _ ports.CategoryRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllCategory(ctx context.Context) ([]domain.Category, error) {
	var models []CategoryModel

	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	colors := make([]domain.Category, 0, len(models))

	for _, model := range models {
		colors = append(colors, domain.Category{
			ID:        model.ID,
			Name:      model.Name,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		})
	}

	return colors, nil
}
