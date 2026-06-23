package ports

import (
	"auth-service/internal/domain"
	"context"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string, userType domain.UserType) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}
