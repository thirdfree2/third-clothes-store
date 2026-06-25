package domain

import "errors"

var (
	ErrCustomerProfileNotFound = errors.New("customer profile not found")
	ErrCustomerAddressNotFound = errors.New("customer address not found")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrCartNotFound            = errors.New("cart not found")
	ErrCartItemNotFound        = errors.New("cart item not found")
	ErrInvalidCartItem         = errors.New("invalid cart item")
	ErrInvalidPurchaseItem     = errors.New("invalid purchase item")
)
