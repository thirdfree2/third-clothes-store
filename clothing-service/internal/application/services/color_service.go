package services

import (
	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"
	"context"
)

type ColorService struct {
	repo ports.ColorRepository
}

func NewColorService(
	repo ports.ColorRepository,
) *ColorService {
	return &ColorService{
		repo: repo,
	}
}

func (s *ColorService) List(
	ctx context.Context) ([]domain.Color, error) {
	return s.repo.GetAllColor(ctx)
}
