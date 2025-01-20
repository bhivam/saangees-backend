package data

import (
	"time"
)

const ScopeAuthentication = "authentication"

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenStore interface {
	CreateToken(userID int64, ttl time.Duration, scope string) (*Token, error)
	InsertToken(token *Token) error
	DeleteAllForUser(scope string, userID int64) error
  GetToken(scope string, plaintext string) (*Token, error)
}
