package data

import (
	"time"

	"github.com/bhivam/saangees-backend/util"
)

var AnonymousUser = &User{}

type User struct {
	ID          int64     `gorm:"primaryKey"                     json:"id"`
	FirstName   string    `gorm:"type:text;not null"             json:"first_name"`
	LastName    string    `gorm:"type:text;not null"             json:"last_name"`
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

// TODO: do password checks separately upon sign up
func ValidateUser(v *util.Validator, user *User) {
	v.Check(util.Nonempty(user.FirstName), "first_name", "Must exist")
	v.Check(util.MaxLen(user.FirstName, 50), "first_name", "Maximum length is 50")
	v.Check(util.Nonempty(user.LastName), "last_name", "Must exist")
	v.Check(util.MaxLen(user.LastName, 50), "last_name", "Maximum length is 50")
	v.Check(util.Nonempty(user.PhoneNumber), "phone_number", "Must exist")

	v.Check(
		util.Matches(user.PhoneNumber, util.PhoneRX),
		"phone_number",
		"Invalid phone number format",
	)
}
