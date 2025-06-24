package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) *JwtMaker {
	return &JwtMaker{secretKey: secretKey}
}

func (maker *JwtMaker) GenerateToken(email string, duration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.secretKey))
}
func (maker *JwtMaker) VerifyToken(tokenStr string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil // email or user id depending on what you set
}
