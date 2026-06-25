package services

import (
	"context"
	"strings"

	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"
)

type ClothesService struct {
	repo          ports.ClothesRepository
	imageRepo     ports.ClothesImageRepository
	changeLogRepo ports.ClothesChangeLogRepository
}

func NewClothesService(
	repo ports.ClothesRepository,
	imageRepo ports.ClothesImageRepository,
	changeLogRepo ports.ClothesChangeLogRepository,
) *ClothesService {
	return &ClothesService{
		repo:          repo,
		imageRepo:     imageRepo,
		changeLogRepo: changeLogRepo,
	}
}

func (s *ClothesService) Create(
	ctx context.Context,
	name string,
	price float64,
	colorID *int64,
	categoryIDs []int64,
) (*domain.Clothes, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, domain.ErrInvalidClothesName
	}

	if err := validatePrice(price); err != nil {
		return nil, err
	}

	if err := validateColorID(colorID); err != nil {
		return nil, err
	}

	if err := validateCategoryIDs(categoryIDs); err != nil {
		return nil, err
	}

	clothes := &domain.Clothes{
		Name:    name,
		Price:   price,
		ColorID: colorID,
	}

	if err := s.repo.Create(ctx, clothes, categoryIDs); err != nil {
		return nil, err
	}

	log := &domain.ClothesChangeLog{
		ClothesID:  &clothes.ID,
		Action:     domain.ClothesChangeActionCreate,
		BeforeData: nil,
		AfterData:  clothesToLogData(clothes),
	}

	if err := s.changeLogRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	return clothes, nil
}

func (s *ClothesService) Update(
	ctx context.Context,
	id int64,
	name string,
	price float64,
	colorID *int64,
	categoryIDs []int64,
) (*domain.Clothes, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidClothesID
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidClothesName
	}

	if err := validatePrice(price); err != nil {
		return nil, err
	}

	if err := validateColorID(colorID); err != nil {
		return nil, err
	}

	if err := validateCategoryIDs(categoryIDs); err != nil {
		return nil, err
	}

	before, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	clothes := &domain.Clothes{
		ID:      id,
		Name:    name,
		Price:   price,
		ColorID: colorID,
	}

	if err := s.repo.Update(ctx, clothes, categoryIDs); err != nil {
		return nil, err
	}

	log := &domain.ClothesChangeLog{
		ClothesID:  &id,
		Action:     domain.ClothesChangeActionUpdate,
		BeforeData: clothesToLogData(before),
		AfterData:  clothesToLogData(clothes),
	}

	if err := s.changeLogRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	return clothes, nil
}

func (s *ClothesService) FindByID(ctx context.Context, id int64) (*domain.Clothes, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidClothesID
	}

	return s.repo.FindByID(ctx, id)
}

func (s *ClothesService) List(
	ctx context.Context,
	pagination ports.Pagination,
	filter ports.ClothesListFilter,
) ([]domain.Clothes, int64, error) {
	if err := validateCategoryID(filter.CategoryID); err != nil {
		return nil, 0, err
	}

	return s.repo.List(ctx, pagination, filter)
}

func (s *ClothesService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidClothesID
	}

	before, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	images, err := s.imageRepo.ListByClothesID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.imageRepo.DeleteByClothesID(ctx, id); err != nil {
		return err
	}

	beforeData := clothesToLogData(before)
	beforeData["images"] = clothesImagesToLogData(images)

	for _, image := range images {
		image := image

		log := &domain.ClothesChangeLog{
			ClothesID: &id,
			Action:    domain.ClothesChangeActionImageDelete,
			BeforeData: map[string]any{
				"image":  clothesImageToLogData(&image),
				"reason": "deleted because clothes was deleted",
			},
			AfterData: nil,
		}

		if err := s.changeLogRepo.Create(ctx, log); err != nil {
			return err
		}
	}

	return nil
}

func validatePrice(price float64) error {
	if price < 0 {
		return domain.ErrInvalidClothesPrice
	}

	return nil
}

func validateColorID(colorID *int64) error {
	if colorID == nil {
		return nil
	}

	if *colorID <= 0 {
		return domain.ErrInvalidColorID
	}

	return nil
}

func validateCategoryIDs(categoryIDs []int64) error {
	for _, categoryID := range categoryIDs {
		if categoryID <= 0 {
			return domain.ErrInvalidCategoryID
		}
	}

	return nil
}

func validateCategoryID(categoryID *int64) error {
	if categoryID == nil {
		return nil
	}

	if *categoryID <= 0 {
		return domain.ErrInvalidCategoryID
	}

	return nil
}
