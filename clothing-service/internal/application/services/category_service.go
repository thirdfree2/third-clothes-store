package services

import (
	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"
	"context"
)

type CategoryService struct {
	repo ports.CategoryRepository
}

func NewCategoryService(
	repo ports.CategoryRepository,
) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) List(
	ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetAllCategory(ctx)
}
