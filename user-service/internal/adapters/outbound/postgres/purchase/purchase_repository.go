package postgrespurchase

import (
	"context"
	"strings"
	"user-service/internal/application/ports"
	"user-service/internal/domain"

	"gorm.io/gorm"
)

var _ ports.PurchaseRepository = (*Repository)(nil)

type Repository struct {
	db             *gorm.DB
	minIOPublicURL string
}

func NewRepository(db *gorm.DB, minIOPublicURL string) *Repository {
	return &Repository{
		db:             db,
		minIOPublicURL: strings.TrimRight(minIOPublicURL, "/"),
	}
}

func (r *Repository) Create(ctx context.Context, purchase domain.Purchase) (*domain.Purchase, error) {
	model := PurchaseModel{
		UserID:      purchase.UserID,
		TotalAmount: purchase.TotalAmount,
		Items:       make([]PurchaseItemModel, 0, len(purchase.Items)),
	}

	for _, item := range purchase.Items {
		model.Items = append(model.Items, PurchaseItemModel{
			ClothesID:       item.ClothesID,
			ProductName:     item.ProductName,
			ProductImageURL: item.ProductImageURL,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			Size:            item.Size,
		})
	}

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
			return err
		}

		return tx.WithContext(ctx).
			Where("purchase_id = ?", model.ID).
			Order("id ASC").
			Find(&model.Items).Error
	}); err != nil {
		return nil, err
	}

	purchaseDomain := toDomain(model)
	return &purchaseDomain, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID int64) ([]domain.Purchase, error) {
	var models []PurchaseModel

	if err := r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Where("user_id = ?", userID).
		Order("created_at DESC, id DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	if err := r.fillMissingProductSnapshots(ctx, models); err != nil {
		return nil, err
	}

	purchases := make([]domain.Purchase, 0, len(models))
	for _, model := range models {
		purchases = append(purchases, toDomain(model))
	}

	return purchases, nil
}

type productSnapshotRow struct {
	ClothesID       int64
	ProductName     string
	ProductImageURL *string
}

func (r *Repository) fillMissingProductSnapshots(ctx context.Context, models []PurchaseModel) error {
	clothesIDs := make(map[int64]struct{})

	for _, model := range models {
		for _, item := range model.Items {
			if item.ProductName == "" || item.ProductImageURL == nil {
				clothesIDs[item.ClothesID] = struct{}{}
			}
		}
	}

	if len(clothesIDs) == 0 {
		return nil
	}

	ids := make([]int64, 0, len(clothesIDs))
	for id := range clothesIDs {
		ids = append(ids, id)
	}

	var rows []productSnapshotRow
	if err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				c.id AS clothes_id,
				c.name AS product_name,
				CASE
					WHEN first_image.id IS NULL THEN NULL
					ELSE ? || '/' || first_image.bucket || '/' || first_image.object_key
				END AS product_image_url
			FROM clothes AS c
			LEFT JOIN LATERAL (
				SELECT ci.id, ci.bucket, ci.object_key
				FROM clothes_images AS ci
				WHERE ci.clothes_id = c.id
				ORDER BY ci.id ASC
				LIMIT 1
			) AS first_image ON true
			WHERE c.id IN ?
		`, r.minIOPublicURL, ids).
		Scan(&rows).Error; err != nil {
		return err
	}

	snapshots := make(map[int64]productSnapshotRow, len(rows))
	for _, row := range rows {
		snapshots[row.ClothesID] = row
	}

	for modelIndex := range models {
		for itemIndex := range models[modelIndex].Items {
			item := &models[modelIndex].Items[itemIndex]
			snapshot, ok := snapshots[item.ClothesID]
			if !ok {
				continue
			}

			if item.ProductName == "" {
				item.ProductName = snapshot.ProductName
			}
			if item.ProductImageURL == nil {
				item.ProductImageURL = snapshot.ProductImageURL
			}
		}
	}

	return nil
}
