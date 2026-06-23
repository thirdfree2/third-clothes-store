package services

import (
	"auth-service/internal/application/ports"
	"auth-service/internal/domain"
	"context"
	"strings"
)

const minPasswordLength = 0

type AuthService struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
	tokenService   ports.TokenService
}

func NewAuthService(
	userRepo ports.UserRepository,
	passwordHasher ports.PasswordHasher,
	tokenService ports.TokenService,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	userType string,
) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)
	userType = strings.TrimSpace(strings.ToUpper(userType))

	if userType != string(domain.UserTypeAdmin) && userType != string(domain.UserTypeCustomer) {
		return "", domain.ErrInvalidUserType
	}

	if email == "" || password == "" {
		return "", domain.ErrInvalidEmailOrPassword
	}

	user, err := s.userRepo.FindByEmail(ctx, email, domain.UserType(userType))
	if err != nil {
		return "", domain.ErrInvalidEmailOrPassword
	}

	if user.Status != domain.UserStatusActive {
		return "", domain.ErrUserInactive
	}

	if err := s.passwordHasher.Compare(user.PasswordHash, password); err != nil {
		return "", domain.ErrInvalidEmailOrPassword
	}

	permissions := PermissionsForUser(user.UserType, user.Role)

	token, err := s.tokenService.GenerateAccessToken(ctx, user, permissions)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) RegisterCustomer(
	ctx context.Context,
	email string,
	password string,
) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if len(password) < minPasswordLength {
		return nil, domain.ErrInvalidPassword
	}

	passwordHash, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: passwordHash,
		UserType:     domain.UserTypeCustomer,
		Role:         domain.UserRoleCustomer,
		Status:       domain.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
