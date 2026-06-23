package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"

	"github.com/google/uuid"
)

type ClothesImageService struct {
	clothesRepo   ports.ClothesRepository
	imageRepo     ports.ClothesImageRepository
	storage       ports.ObjectStorage
	changeLogRepo ports.ClothesChangeLogRepository
	bucket        string
}

func NewClothesImageService(
	clothesRepo ports.ClothesRepository,
	imageRepo ports.ClothesImageRepository,
	storage ports.ObjectStorage,
	changeLogRepo ports.ClothesChangeLogRepository,
	bucket string,
) *ClothesImageService {
	return &ClothesImageService{
		clothesRepo:   clothesRepo,
		imageRepo:     imageRepo,
		storage:       storage,
		changeLogRepo: changeLogRepo,
		bucket:        bucket,
	}
}

func (s *ClothesImageService) Upload(
	ctx context.Context,
	clothesID int64,
	filename string,
	contentType string,
	size int64,
	reader io.Reader,
) (*domain.ClothesImage, error) {
	if clothesID <= 0 {
		return nil, domain.ErrInvalidClothesID
	}

	// เช็กว่า clothes มีอยู่จริง
	if _, err := s.clothesRepo.FindByID(ctx, clothesID); err != nil {
		return nil, err
	}

	if !isAllowedImageContentType(contentType) {
		return nil, fmt.Errorf("unsupported image content type: %s", contentType)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	objectKey := fmt.Sprintf(
		"clothes/%d/%s%s",
		clothesID,
		uuid.NewString(),
		ext,
	)

	if err := s.storage.Upload(ctx, ports.UploadObjectInput{
		Bucket:      s.bucket,
		ObjectKey:   objectKey,
		Reader:      reader,
		Size:        size,
		ContentType: contentType,
	}); err != nil {
		return nil, err
	}

	image := &domain.ClothesImage{
		ClothesID:        clothesID,
		Bucket:           s.bucket,
		ObjectKey:        objectKey,
		OriginalFilename: &filename,
		ContentType:      contentType,
		SizeBytes:        size,
		CreatedAt:        time.Now(),
	}

	if err := s.imageRepo.Create(ctx, image); err != nil {
		return nil, err
	}

	log := &domain.ClothesChangeLog{
		ClothesID:  &clothesID,
		Action:     domain.ClothesChangeActionImageUpload,
		BeforeData: nil,
		AfterData: map[string]any{
			"image": clothesImageToLogData(image),
		},
	}

	if err := s.changeLogRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	return image, nil
}

func (s *ClothesImageService) Delete(
	ctx context.Context,
	clothesID int64,
	imageID int64,
) error {
	if clothesID <= 0 {
		return domain.ErrInvalidClothesID
	}

	if imageID <= 0 {
		return domain.ErrInvalidClothesImageID
	}

	image, err := s.imageRepo.FindByID(ctx, imageID)
	if err != nil {
		return err
	}

	if image.ClothesID != clothesID {
		return domain.ErrClothesImageNotFound
	}

	if err := s.imageRepo.Delete(ctx, imageID); err != nil {
		return err
	}

	log := &domain.ClothesChangeLog{
		ClothesID: &clothesID,
		Action:    domain.ClothesChangeActionImageDelete,
		BeforeData: map[string]any{
			"image": clothesImageToLogData(image),
		},
		AfterData: nil,
	}

	if err := s.changeLogRepo.Create(ctx, log); err != nil {
		return err
	}

	return nil
}

func isAllowedImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}
