package postgresuser

import (
	"time"

	"auth-service/internal/domain"
)

type UserModel struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Email        string    `gorm:"column:email;size:255;not null;uniqueIndex"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	UserType     string    `gorm:"column:user_type;size:20;not null"`
	Status       string    `gorm:"column:status;size:20;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

type AdminProfileModel struct {
	UserID    int64     `gorm:"column:user_id;primaryKey"`
	Role      string    `gorm:"column:role;size:50;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (AdminProfileModel) TableName() string {
	return "admin_profiles"
}

type userRecord struct {
	ID           int64
	Email        string
	PasswordHash string
	UserType     string
	Role         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func toDomain(model userRecord) domain.User {
	return domain.User{
		ID:           model.ID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		UserType:     domain.UserType(model.UserType),
		Role:         domain.UserRole(model.Role),
		Status:       domain.UserStatus(model.Status),
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}
