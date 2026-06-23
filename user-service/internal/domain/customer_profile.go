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
