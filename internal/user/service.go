package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type service interface {
	hashPassword(password string) (string, error)
	checkPassword(hashedPassword, password string) error
}

type userService struct{}

func newUserService() service {
	return &userService{}
}

func (s *userService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *userService) checkPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
