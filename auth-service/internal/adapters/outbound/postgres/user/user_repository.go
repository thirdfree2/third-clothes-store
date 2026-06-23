package postgresuser

import (
	"context"
	"errors"
	"strings"

	"auth-service/internal/application/ports"
	"auth-service/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var _ ports.UserRepository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByEmail(ctx context.Context, email string, userType domain.UserType) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	normalizedUserType := domain.UserType(strings.TrimSpace(strings.ToUpper(string(userType))))

	var model userRecord

	if err := r.db.WithContext(ctx).
		Table("users AS u").
		Select(`
			u.id,
			u.email,
			u.password_hash,
			u.user_type,
			u.status,
			u.created_at,
			u.updated_at,
			CASE
				WHEN u.user_type = ? THEN ap.role
				WHEN u.user_type = ? THEN ?
				ELSE u.user_type
			END AS role
		`, domain.UserTypeAdmin, domain.UserTypeCustomer, domain.UserRoleCustomer).
		Joins("LEFT JOIN admin_profiles AS ap ON ap.user_id = u.id").
		Where("LOWER(u.email) = ?", email).
		Where("u.user_type = ?", normalizedUserType).
		Take(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidEmailOrPassword
		}

		return nil, err
	}

	admin := toDomain(model)
	return &admin, nil
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	model := UserModel{
		Email:        strings.TrimSpace(strings.ToLower(user.Email)),
		PasswordHash: user.PasswordHash,
		UserType:     string(user.UserType),
		Status:       string(user.Status),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrEmailDuplicated
		}

		return err
	}

	user.ID = model.ID
	user.Email = model.Email
	user.UserType = domain.UserType(model.UserType)
	user.Status = domain.UserStatus(model.Status)
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}
