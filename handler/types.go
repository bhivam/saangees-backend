package handler

import (
	"time"

	"github.com/bhivam/saangees-backend/data"
)

type UserRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserResponseBody struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type ListUsersResponseBody []UserResponseBody

type LoginUserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginUserResponseBody struct {
	SessionID             string           `json:"session_id"`
	AccessToken           string           `json:"access_token"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	User                  UserResponseBody `json:"user"`
}

type RefreshTokenResponseBody struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func toDataUser(userRequestBody UserRequestBody) *data.User {
	return &data.User{
		Name:     userRequestBody.Name,
		Email:    userRequestBody.Email,
		Password: userRequestBody.Password,
		IsAdmin:  userRequestBody.IsAdmin,
	}
}

func toUserRes(user *data.User) *UserResponseBody {
	return &UserResponseBody{
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}
}
