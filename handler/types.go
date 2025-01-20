package handler

import (
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

type CreateTokenRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenResponseBody struct {
	User  UserResponseBody `json:"user"`
	Token *data.Token      `json:"token"`
}

func toDataUser(userRequestBody UserRequestBody) *data.User {
	return &data.User{
		Name:    userRequestBody.Name,
		Email:   userRequestBody.Email,
		Hash:    userRequestBody.Password, // prehashed
		IsAdmin: userRequestBody.IsAdmin,
	}
}

func toUserRes(user *data.User) *UserResponseBody {
	return &UserResponseBody{
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}
}
