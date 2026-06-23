package services

import (
	"context"
	"math"
	"strings"
	"user-service/internal/application/ports"
	"user-service/internal/domain"
)

type PurchaseService struct {
	purchaseRepo  ports.PurchaseRepository
	clothesClient ports.ClothesClient
}

func NewPurchaseService(purchaseRepo ports.PurchaseRepository, clothesClient ports.ClothesClient) *PurchaseService {
	return &PurchaseService{
		purchaseRepo:  purchaseRepo,
		clothesClient: clothesClient,
	}
}

func (s *PurchaseService) CreatePurchase(ctx context.Context, userID int64, items []domain.PurchaseItem) (*domain.Purchase, error) {
	if userID <= 0 || len(items) == 0 {
		return nil, domain.ErrInvalidPurchaseItem
	}

	purchaseItems := make([]domain.PurchaseItem, 0, len(items))
	clothesByID := make(map[int64]*domain.Clothes)
	var totalAmount float64

	for _, item := range items {
		if item.ClothesID <= 0 || item.Quantity <= 0 {
			return nil, domain.ErrInvalidPurchaseItem
		}

		clothes, ok := clothesByID[item.ClothesID]
		if !ok {
			var err error
			clothes, err = s.clothesClient.GetByID(ctx, item.ClothesID)
			if err != nil {
				return nil, err
			}
			clothesByID[item.ClothesID] = clothes
		}

		unitPrice := roundMoney(clothes.Price)
		totalAmount += unitPrice * float64(item.Quantity)

		purchaseItems = append(purchaseItems, domain.PurchaseItem{
			ClothesID:       item.ClothesID,
			ProductName:     strings.TrimSpace(clothes.Name),
			ProductImageURL: firstImageURL(clothes.Images),
			Quantity:        item.Quantity,
			UnitPrice:       unitPrice,
			Size:            strings.TrimSpace(item.Size),
		})
	}

	return s.purchaseRepo.Create(ctx, domain.Purchase{
		UserID:      userID,
		TotalAmount: roundMoney(totalAmount),
		Items:       purchaseItems,
	})
}

func (s *PurchaseService) ListPurchases(ctx context.Context, userID int64) ([]domain.Purchase, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidPurchaseItem
	}

	return s.purchaseRepo.FindByUserID(ctx, userID)
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func firstImageURL(images []domain.ClothesImage) *string {
	if len(images) == 0 || strings.TrimSpace(images[0].ImageURL) == "" {
		return nil
	}

	imageURL := strings.TrimSpace(images[0].ImageURL)
	return &imageURL
}
