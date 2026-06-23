package domain

import "errors"

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrUserInactive           = errors.New("user inactive")
	ErrInvalidAdminRole       = errors.New("invalid admin role")
	ErrInvalidUserType        = errors.New("invalid user type")
	ErrEmailDuplicated        = errors.New("email already exists")
	ErrInvalidPassword        = errors.New("invalid password")
)
