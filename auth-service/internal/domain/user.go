package domain

import "time"

type UserType string
type UserRole string
type UserStatus string

const (
	UserRoleAdminViewer UserRole = "ADMIN_VIEWER"
	UserRoleAdminEditor UserRole = "ADMIN_EDITOR"
	UserRoleAdminOwner  UserRole = "ADMIN_OWNER"
	UserRoleCustomer    UserRole = "CUSTOMER"

	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusInactive UserStatus = "INACTIVE"
	UserStatusBanned   UserStatus = "BANNED"

	UserTypeAdmin    UserType = "ADMIN"
	UserTypeCustomer UserType = "CUSTOMER"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	UserType     UserType
	Role         UserRole
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
