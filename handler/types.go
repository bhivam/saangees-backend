package handler

import (
	"time"

	"github.com/bhivam/saangees-backend/data"
)

type UserRequestBody struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserResponseBody struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
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
		FirstName:   userRequestBody.FirstName,
		LastName:    userRequestBody.LastName,
		PhoneNumber: userRequestBody.PhoneNumber,
		Hash:        userRequestBody.Password, // prehashed
		IsAdmin:     false,
	}
}

type CreateItemRequest struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	BasePrice    float64            `json:"base_price"`
	Date         time.Time          `json:"date"`
	SpiceOptions []data.SpiceOption `json:"spice_options"`
}

type UpdateItemRequest struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	BasePrice    float64            `json:"base_price"`
	Date         time.Time          `json:"date"`
	SpiceOptions []data.SpiceOption `json:"spice_options"`
	Published    bool               `json:"published"`
}

type ItemResponse struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Date         time.Time          `json:"date"`
	Description  string             `json:"description"`
	BasePrice    float64            `json:"base_price"`
	SpiceOptions []data.SpiceOption `json:"spice_options"`
	Published    bool               `json:"published"`
}

type ListItemsResponse []ItemResponse

func toUserRes(user *data.User) *UserResponseBody {
	return &UserResponseBody{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsAdmin:     user.IsAdmin,
	}
}

func toItemResponse(item *data.Item) *ItemResponse {
	res := &ItemResponse{
		ID:           item.ID,
		Name:         item.Name,
		Date:         item.Date,
		Description:  item.Description,
		BasePrice:    item.BasePrice,
		SpiceOptions: item.SpiceOptions,
		Published:    item.Published,
	}

	return res
}
