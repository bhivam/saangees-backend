package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID      int64  `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewUserClaims(
	id int64,
	email string,
	isAdmin bool,
	duration time.Duration,
) (*UserClaims, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("Error creating UUID")
	}

	return &UserClaims{
		ID:      id,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}, nil
}
