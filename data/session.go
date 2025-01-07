package data

import "time"

type Session struct {
	ID           string    `json:"id"            db:"id"`
	UserEmail    string    `json:"user_email"    db:"user_email"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"    db:"is_revoked"`
	ExpiresAt    time.Time `json:"expires_at"    db:"expires_at"`
	CreatedAt    time.Time `json:"created_at"    db:"created_at"`
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetSession(id string) (*Session, error)
	RevokedSession(id string) error
	DeleteSession(id string) error
}

// implement the SessionStore interface for a in memory array of sessions
