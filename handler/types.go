package handler

import (
	"time"

	"github.com/bhivam/saangees-backend/data"
)

// ========== USER ==========

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

func toDataUser(userRequestBody UserRequestBody) *data.User {
	return &data.User{
		FirstName:   userRequestBody.FirstName,
		LastName:    userRequestBody.LastName,
		PhoneNumber: userRequestBody.PhoneNumber,
		Hash:        userRequestBody.Password, // prehashed
		IsAdmin:     false,
	}
}

func toUserRes(user *data.User) *UserResponseBody {
	return &UserResponseBody{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsAdmin:     user.IsAdmin,
	}
}

// ========== TOKEN ==========

type CreateTokenRequestBody struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// ========== MODIFER ==========

// NOTE: thse exist purely to strip IDs from the db definition--ID unknown in create req
type ModifierOption struct {
	Name          string  `json:"name"`
	PriceModifier float64 `json:"price_modifier"`
}

type ModifierOptions []ModifierOption

type ModifierCategory struct {
	Name            string          `json:"name"`
	Min             int8            `json:"min"`
	Max             int8            `json:"max"`
	ModifierOptions ModifierOptions `json:"modifier_options"`
}

type ModiferCategories []ModifierCategory

// ========== ITEM ==========

type CreateItemRequest struct {
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	BasePrice          float64           `json:"base_price"`
	Date               time.Time         `json:"date"`
	Quantity           uint8             `json:"quantity"`
	Unit               string            `json:"unit"`
	ModifierCategories ModiferCategories `json:"modifier_categories"`
}

type ListItemsResponse []data.Item

func toDataItem(req CreateItemRequest) data.Item {
	modifierCategories := []data.ModifierCategory{}

	if len(req.ModifierCategories) > 0 {
		for _, mc := range req.ModifierCategories {
			modifierOptions := []data.ModifierOption{}

			for _, mo := range mc.ModifierOptions {
				modifierOptions = append(modifierOptions, data.ModifierOption{
					Name:          mo.Name,
					PriceModifier: mo.PriceModifier,
				})
			}

			modifierCategories = append(modifierCategories, data.ModifierCategory{
				Name:            mc.Name,
				Min:             mc.Min,
				Max:             mc.Max,
				ModifierOptions: modifierOptions,
			})
		}
	}

	return data.Item{
		Name:               req.Name,
		Description:        req.Description,
		BasePrice:          req.BasePrice,
		Date:               req.Date,
		Quantity:           req.Quantity,
		Unit:               req.Unit,
		ModifierCategories: modifierCategories,
	}
}
