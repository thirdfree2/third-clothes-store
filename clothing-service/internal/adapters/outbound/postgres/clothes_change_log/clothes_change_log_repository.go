package postgresclotheschangelog

import (
	"clothing-service/internal/application/ports"
	"clothing-service/internal/domain"
	"context"
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var _ ports.ClothesChangeLogRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, log *domain.ClothesChangeLog) error {
	if log == nil {
		return errors.New("clothes change log is nil")
	}

	beforeData, err := marshalJSONB(log.BeforeData)
	if err != nil {
		return err
	}

	afterData, err := marshalJSONB(log.AfterData)
	if err != nil {
		return err
	}

	model := &clothesChangeLogModel{
		ClothesID:  log.ClothesID,
		Action:     string(log.Action),
		BeforeData: beforeData,
		AfterData:  afterData,
		ChangedBy:  log.ChangedBy,
		Note:       log.Note,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	log.ID = model.ID
	log.CreatedAt = model.CreatedAt

	return nil
}

func marshalJSONB(value map[string]any) (datatypes.JSON, error) {
	if value == nil {
		return nil, nil
	}

	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(b), nil
}
