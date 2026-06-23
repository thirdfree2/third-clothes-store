package services

import (
	"context"
	"user-service/internal/application/ports"
	"user-service/internal/domain"
)

type CustomerService struct {
	customerProfileRepo ports.CustomerProfileRepository
}

func NewCustomerService(customerProfileRepo ports.CustomerProfileRepository) *CustomerService {
	return &CustomerService{
		customerProfileRepo: customerProfileRepo,
	}
}

func (s *CustomerService) GetProfile(ctx context.Context, userID int64) (*domain.CustomerProfile, error) {
	return s.customerProfileRepo.FindByUserID(ctx, userID)
}

func (s *CustomerService) UpdateProfile(ctx context.Context, profile *domain.CustomerProfile) error {
	return s.customerProfileRepo.Upsert(ctx, profile)
}
