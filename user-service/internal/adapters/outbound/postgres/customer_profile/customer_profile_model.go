package postgrescustomerprofile

import (
	"time"
	"user-service/internal/domain"
)

type CustomerProfileModel struct {
	UserID         int64      `gorm:"column:user_id;primaryKey"`
	FirstName      *string    `gorm:"column:first_name;size:100"`
	LastName       *string    `gorm:"column:last_name;size:100"`
	Phone          *string    `gorm:"column:phone;size:30"`
	DateOfBirth    *time.Time `gorm:"column:date_of_birth"`
	MarketingOptIn bool       `gorm:"column:marketing_opt_in;not null"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (CustomerProfileModel) TableName() string {
	return "customer_profiles"
}

func toModel(profile *domain.CustomerProfile) CustomerProfileModel {
	return CustomerProfileModel{
		UserID:         profile.UserID,
		FirstName:      profile.FirstName,
		LastName:       profile.LastName,
		Phone:          profile.Phone,
		DateOfBirth:    profile.DateOfBirth,
		MarketingOptIn: profile.MarketingOptIn,
		CreatedAt:      profile.CreatedAt,
		UpdatedAt:      profile.UpdatedAt,
	}
}

func toDomain(model CustomerProfileModel) domain.CustomerProfile {
	return domain.CustomerProfile{
		UserID:         model.UserID,
		FirstName:      model.FirstName,
		LastName:       model.LastName,
		Phone:          model.Phone,
		DateOfBirth:    model.DateOfBirth,
		MarketingOptIn: model.MarketingOptIn,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}
