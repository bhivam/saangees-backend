package handler

import (
	"time"

	"github.com/bhivam/saangees-backend/data"
)

type UserRequestBody struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	IsAdmin     bool   `json:"is_admin"`
}

type UserResponseBody struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	IsAdmin     bool   `json:"is_admin"`
}

type ListUsersResponseBody []UserResponseBody

type CreateTokenRequestBody struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func toDataUser(userRequestBody UserRequestBody) *data.User {
	return &data.User{
		Name:        userRequestBody.Name,
		PhoneNumber: userRequestBody.PhoneNumber,
		Hash:        userRequestBody.Password, // prehashed
		IsAdmin:     userRequestBody.IsAdmin,
	}
}

type CreateItemRequest struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	BasePrice    float64            `json:"base_price"`
	Date         time.Time          `json:"date"`
	SizeOptions  []data.SizeOption  `json:"size_options"`
	SpiceOptions []data.SpiceOption `json:"spice_options"`
}

type UpdateItemRequest struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	BasePrice    float64            `json:"base_price"`
	Date         time.Time          `json:"date"`
	SizeOptions  []data.SizeOption  `json:"size_options"`
	SpiceOptions []data.SpiceOption `json:"spice_options"`
} 

type ItemResponse struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Date         time.Time          `json:"date"`
	Description  string             `json:"description"`
	BasePrice    float64            `json:"base_price"`
	SizeOptions  []data.SizeOption  `json:"size_options"`
	SpiceOptions []data.SpiceOption `json:"spice_options"`
}

type ListItemsResponse []ItemResponse

func toUserRes(user *data.User) *UserResponseBody {
	return &UserResponseBody{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		IsAdmin:     user.IsAdmin,
	}
}

func toItemResponse(item *data.Item) *ItemResponse {
	return &ItemResponse{
		ID:   item.ID,
		Name: item.Name,
		Date: item.Date,
    Description: item.Description,
    BasePrice: item.BasePrice,
    SizeOptions: item.SizeOptions,
    SpiceOptions: item.SpiceOptions,
	}
}
