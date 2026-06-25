package services

import (
	"context"
	"user-service/internal/application/ports"
	"user-service/internal/domain"
)

type CustomerService struct {
	customerProfileRepo ports.CustomerProfileRepository
	customerAddressRepo ports.CustomerAddressRepository
}

func NewCustomerService(
	customerProfileRepo ports.CustomerProfileRepository,
	customerAddressRepo ports.CustomerAddressRepository,
) *CustomerService {
	return &CustomerService{
		customerProfileRepo: customerProfileRepo,
		customerAddressRepo: customerAddressRepo,
	}
}

func (s *CustomerService) GetProfile(ctx context.Context, userID int64) (*domain.CustomerProfile, error) {
	return s.customerProfileRepo.FindByUserID(ctx, userID)
}

func (s *CustomerService) UpdateProfile(ctx context.Context, profile *domain.CustomerProfile) error {
	return s.customerProfileRepo.Upsert(ctx, profile)
}

func (s *CustomerService) GetDefaultAddress(ctx context.Context, userID int64) (*domain.CustomerAddress, error) {
	return s.customerAddressRepo.FindDefaultByUserID(ctx, userID)
}

func (s *CustomerService) UpdateDefaultAddress(ctx context.Context, address *domain.CustomerAddress) error {
	address.IsDefault = true
	if address.CountryCode == "" {
		address.CountryCode = "TH"
	}

	return s.customerAddressRepo.UpsertDefault(ctx, address)
}
