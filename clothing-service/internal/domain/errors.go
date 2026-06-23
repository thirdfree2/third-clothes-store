package domain

import "errors"

var (
	ErrClothesDuplicated        = errors.New("clothes already exists")
	ErrClothesNotFound          = errors.New("clothes not found")
	ErrInvalidClothesID         = errors.New("invalid clothes id")
	ErrInvalidClothesName       = errors.New("invalid clothes name")
	ErrInvalidClothesPrice      = errors.New("invalid clothes price")
	ErrInvalidColorID           = errors.New("invalid color id")
	ErrInvalidCategoryID        = errors.New("invalid category id")
	ErrClothesImageNotFound     = errors.New("clothes image not found")
	ErrInvalidClothesImageID    = errors.New("invalid clothes image id")
	ErrClothesNameAlreadyExists = errors.New("clothes name already exists")
)
