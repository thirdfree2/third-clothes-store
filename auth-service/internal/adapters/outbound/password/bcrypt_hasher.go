// internal/adapters/outbound/password/bcrypt_hasher.go

package password

import (
	"auth-service/internal/application/ports"

	"golang.org/x/crypto/bcrypt"
)

var _ ports.PasswordHasher = (*BcryptHasher)(nil)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *BcryptHasher) Compare(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
}
