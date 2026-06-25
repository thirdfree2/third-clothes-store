package domain

import "time"

type CustomerProfile struct {
	UserID         int64
	FirstName      *string
	LastName       *string
	Phone          *string
	DateOfBirth    *time.Time
	MarketingOptIn bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CustomerAddress struct {
	ID            int64
	UserID        int64
	RecipientName string
	Phone         string
	AddressLine1  string
	AddressLine2  *string
	Subdistrict   *string
	District      string
	Province      string
	PostalCode    string
	CountryCode   string
	IsDefault     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
