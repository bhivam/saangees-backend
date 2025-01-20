package data

import (
	"time"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"         db:"id"`
	Name      string    `json:"name"       db:"name"`
	Email     string    `json:"email"      db:"email"`
	Hash      string    `json:"hash"       db:"hash"`
	IsAdmin   bool      `json:"is_admin"   db:"is_admin"` // TODO do role management better later
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserStore interface {
	CreateUser(user *User) (*User, error)
	GetUser(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
  GetByToken(token string, scope string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
	ListUsers() ([]*User, error)
}

func (user *User) IsAnonymous() bool {
  return user == AnonymousUser 
}

// implement the UserStore interface for a in memory array of users
