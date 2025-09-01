package utils

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

type passwordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) PasswordHasher {
	return &passwordHasher{cost: cost}
}

// Hash implements PasswordHasher.
func (p *passwordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify implements PasswordHasher.
func (p *passwordHasher) Verify(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

