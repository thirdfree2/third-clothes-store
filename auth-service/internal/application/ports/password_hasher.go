package ports

type PasswordHasher interface {
	Hash(plainPassword string) (string, error)
	Compare(hashedPassword string, plainPassword string) error
}
