package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (jwtMaker *JWTMaker) GenerateToken(
	id int64,
	email string,
	isAdmin bool,
	duration time.Duration,
) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, isAdmin, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtMaker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("Error signing token: %v", err)
	}

	return signedString, claims, nil
}

func (jwtMaker *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Invalid token signing method")
			}

			return []byte(jwtMaker.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Error casting claims: %v", err)
	}

	return claims, nil
}
