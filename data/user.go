package data

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"         db:"id"`
	Name      string    `json:"name"       db:"name"`
	Email     string    `json:"email"      db:"email"`
	Password  string    `json:"password"   db:"password"`
	IsAdmin   bool      `json:"is_admin"   db:"is_admin"` // do role management better later
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserStore interface {
  CreateUser(user *User) (*User, error)
  GetUser(email string) (*User, error)
  UpdateUser(user *User) error
  DeleteUser(id int64) error
  ListUsers() ([]*User, error)
}

// implement the UserStore interface for a in memory array of users
