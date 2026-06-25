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

type CustomerAddressModel struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64     `gorm:"column:user_id;not null"`
	RecipientName string    `gorm:"column:recipient_name;size:150;not null"`
	Phone         string    `gorm:"column:phone;size:30;not null"`
	AddressLine1  string    `gorm:"column:address_line1;not null"`
	AddressLine2  *string   `gorm:"column:address_line2"`
	Subdistrict   *string   `gorm:"column:subdistrict;size:120"`
	District      string    `gorm:"column:district;size:120;not null"`
	Province      string    `gorm:"column:province;size:120;not null"`
	PostalCode    string    `gorm:"column:postal_code;size:20;not null"`
	CountryCode   string    `gorm:"column:country_code;size:2;not null"`
	IsDefault     bool      `gorm:"column:is_default;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (CustomerAddressModel) TableName() string {
	return "customer_addresses"
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

func addressToModel(address *domain.CustomerAddress) CustomerAddressModel {
	return CustomerAddressModel{
		ID:            address.ID,
		UserID:        address.UserID,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		AddressLine1:  address.AddressLine1,
		AddressLine2:  address.AddressLine2,
		Subdistrict:   address.Subdistrict,
		District:      address.District,
		Province:      address.Province,
		PostalCode:    address.PostalCode,
		CountryCode:   address.CountryCode,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt,
		UpdatedAt:     address.UpdatedAt,
	}
}

func addressToDomain(model CustomerAddressModel) domain.CustomerAddress {
	return domain.CustomerAddress{
		ID:            model.ID,
		UserID:        model.UserID,
		RecipientName: model.RecipientName,
		Phone:         model.Phone,
		AddressLine1:  model.AddressLine1,
		AddressLine2:  model.AddressLine2,
		Subdistrict:   model.Subdistrict,
		District:      model.District,
		Province:      model.Province,
		PostalCode:    model.PostalCode,
		CountryCode:   model.CountryCode,
		IsDefault:     model.IsDefault,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
