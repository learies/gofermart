package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID int, expirationTime time.Time) string
	VerifyToken(tokenString string) (int, error)
}

type jwtService struct{}

func NewJWTService() JWTService {
	return &jwtService{}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func (j *jwtService) GenerateToken(userID int, expirationTime time.Time) string {
	var tokenString string

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ = token.SignedString([]byte("qwerty"))

	return tokenString
}

func (j *jwtService) VerifyToken(tokenString string) (int, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("qwerty"), nil
	})

	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
