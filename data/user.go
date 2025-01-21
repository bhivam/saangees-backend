package data

import (
	"time"
)

var AnonymousUser = &User{}

type User struct {
	ID          int64     `gorm:"primaryKey"                     json:"id"`
	Name        string    `gorm:"type:text;not null"             json:"name"`
	PhoneNumber string    `gorm:"type:text;not null;uniqueIndex" json:"phone_number"`
	Hash        string    `gorm:"type:text;not null"             json:"hash"`
	IsAdmin     bool      `gorm:"default:false"                  json:"is_admin"`
	CreatedAt   time.Time `gorm:"autoCreateTime"                 json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"                 json:"updated_at"`
}

type UserStore interface {
	CreateUser(user *User) (*User, error)
	GetUser(id int64) (*User, error)
	GetByPhoneNumber(phoneNumber string) (*User, error)
	GetByToken(token string, scope string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
	ListUsers() ([]*User, error)
}

func (user *User) IsAnonymous() bool {
	return user == AnonymousUser
}
