package postgrescustomerprofile

import (
	"context"
	"errors"
	"user-service/internal/application/ports"
	"user-service/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CustomerProfileRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUserID(ctx context.Context, userID int64) (*domain.CustomerProfile, error) {
	var model CustomerProfileModel

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Take(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCustomerProfileNotFound
		}

		return nil, err
	}

	profile := toDomain(model)
	return &profile, nil
}

func (r *Repository) Upsert(ctx context.Context, profile *domain.CustomerProfile) error {
	model := toModel(profile)

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"first_name",
				"last_name",
				"phone",
				"date_of_birth",
				"marketing_opt_in",
				"updated_at",
			}),
		}).
		Create(&model).Error
}
