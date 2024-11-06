package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the methods for password hashing and validation.
type AuthService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}

// authService is a concrete implementation of AuthService.
type authService struct{}

// NewAuthService creates a new instance of AuthService.
func NewAuthService() AuthService {
	return &authService{}
}

// HashPassword hashes the given password using bcrypt.
func (s *authService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a hashed password with a plaintext password.
// Returns an error if the passwords do not match.
func (s *authService) VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
