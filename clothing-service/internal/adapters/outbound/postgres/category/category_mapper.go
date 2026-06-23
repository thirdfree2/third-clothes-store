package postgrescategory

import "clothing-service/internal/domain"

func ToCategoryModel(category *domain.Category) CategoryModel {
	if category == nil {
		return CategoryModel{}
	}

	return CategoryModel{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func ToCategoryDomain(model CategoryModel) domain.Category {
	return domain.Category{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
