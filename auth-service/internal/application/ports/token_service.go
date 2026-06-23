package ports

import (
	"auth-service/internal/domain"
	"context"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *domain.User, permissions []string) (string, error)
}
