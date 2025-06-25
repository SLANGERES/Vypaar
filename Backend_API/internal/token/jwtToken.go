package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtMaker struct {
	SecretKey string
}

func NewJwtMaker(secretKey string) *JwtMaker {
	return &JwtMaker{SecretKey: secretKey}
}

type VypaarClaim struct {
	ShopID string `json:"shopID"`
	jwt.RegisteredClaims
}

func (maker *JwtMaker) GenerateToken(email string, shopID string, duration time.Duration) (string, error) {
	claims := VypaarClaim{
		ShopID: shopID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.SecretKey))
}
func (maker *JwtMaker) VerifyToken(tokenString string) (*VypaarClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &VypaarClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*VypaarClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
