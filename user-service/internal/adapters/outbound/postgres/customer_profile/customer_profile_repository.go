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
var _ ports.CustomerAddressRepository = (*Repository)(nil)

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

func (r *Repository) FindDefaultByUserID(ctx context.Context, userID int64) (*domain.CustomerAddress, error) {
	var model CustomerAddressModel

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		Take(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCustomerAddressNotFound
		}

		return nil, err
	}

	address := addressToDomain(model)
	return &address, nil
}

func (r *Repository) UpsertDefault(ctx context.Context, address *domain.CustomerAddress) error {
	model := addressToModel(address)
	model.IsDefault = true

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing CustomerAddressModel
		err := tx.
			Where("user_id = ? AND is_default = ?", model.UserID, true).
			Take(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(&model).Error
		}

		model.ID = existing.ID
		return tx.Model(&existing).Updates(map[string]any{
			"recipient_name": model.RecipientName,
			"phone":          model.Phone,
			"address_line1":  model.AddressLine1,
			"address_line2":  model.AddressLine2,
			"subdistrict":    model.Subdistrict,
			"district":       model.District,
			"province":       model.Province,
			"postal_code":    model.PostalCode,
			"country_code":   model.CountryCode,
			"is_default":     true,
		}).Error
	})
}
