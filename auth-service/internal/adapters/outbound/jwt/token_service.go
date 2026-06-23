package jwtservice

import (
	"context"
	"time"

	"auth-service/internal/application/ports"
	"auth-service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

var _ ports.TokenService = (*TokenService)(nil)

type TokenService struct {
	secret string
	issuer string
	ttl    time.Duration
}

func NewTokenService(secret string, issuer string, ttl time.Duration) *TokenService {
	return &TokenService{
		secret: secret,
		issuer: issuer,
		ttl:    ttl,
	}
}

func (s *TokenService) GenerateAccessToken(
	ctx context.Context,
	user *domain.User,
	permissions []string,
) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":         user.ID,
		"email":       user.Email,
		"user_type":   user.UserType,
		"role":        user.Role,
		"permissions": permissions,
		"iss":         s.issuer,
		"iat":         now.Unix(),
		"exp":         now.Add(s.ttl).Unix(),
		"type":        accessTokenType(user.UserType),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}

func accessTokenType(userType domain.UserType) string {
	switch userType {
	case domain.UserTypeCustomer:
		return "customer_access"
	default:
		return "admin_access"
	}
}
