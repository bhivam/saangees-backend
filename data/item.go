package data

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/bhivam/saangees-backend/util"
)

// NOTE - text vs nvarchar
type Item struct {
	ID                 uint               `json:"id"                  gorm:"primaryKey"`
	Name               string             `json:"name"                gorm:"type:text;not null"`
	Description        string             `json:"description"         gorm:"type:text;not null"`
	BasePrice          float64            `json:"base_price"          gorm:"type:decimal(10,2);not null"`
	Date               time.Time          `json:"date"                gorm:"type:timestamp with time zone;not null"`
	Quantity           uint8              `json:"quantity"            gorm:"not null"`
	Unit               string             `json:"unit"                gorm:"type:text;not null"`
	Published          bool               `json:"published"           gorm:"type:boolean;not null;default:false"`
	ModifierCategories []ModifierCategory `json:"modifier_categories"`
}

type ModifierCategory struct {
	ID              uint             `json:"id"               gorm:"primaryKey"`
	ItemID          uint             `json:"item_id"`
	Name            string           `json:"name"             gorm:"type:text;not null"`
	Min             int8             `json:"min"              gorm:"not null"`
	Max             int8             `json:"max"              gorm:"not null"`
	ModifierOptions []ModifierOption `json:"modifier_options"`
}

type ModifierOption struct {
	ID                 uint    `json:"id"                  gorm:"primaryKey"`
	ModifierCategoryID uint    `json:"modifer_category_id"`
	Name               string  `json:"name"                gorm:"type:text;not null"`
	PriceModifier      float64 `json:"price_modifier"      gorm:"type:decimal(10,2);not null"`
}

type ItemStore interface {
	CreateItem(item *Item) (*Item, error)
	GetItem(id uint) (*Item, error)
	UpdateItem(item *Item) error
	DeleteItem(id uint) error
	ComingWeekItems() ([]*Item, error)
}

func ValidateModifierCategories(v *util.Validator, mc []ModifierCategory) {
	// Name            string
	// Min             int8
	// Max             int8
	// ModifierOptions []ModifierOption

	for i, m := range mc {
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("modifier_categories[%d].", i))

	}
}

func ValidateItem(v *util.Validator, item *Item) {
	// ID                 uint
	// Name               string
	// Description        string
	// BasePrice          float64
	// Date               time.Time
	// Quantity           uint8
	// Unit               string
	// Published          bool
	// ModifierCategories []ModifierCategory

	v.Check(util.Nonempty(item.Name), "name", "Must exist")
	v.Check(util.MaxLen(item.Name, 50), "name", "Maximum length is 50")
	v.Check(util.Nonempty(item.Description), "description", "Must exist")
	v.Check(util.MaxLen(item.Description, 300), "description", "Maximum length is 300")
	v.Check(util.NonNegativeFl(item.BasePrice), "base_price", "Must be nonnegative")
	v.Check(item.Quantity > 0, "quantity", "Must be greater than 0")
	v.Check(util.Nonempty(item.Unit), "unit", "Must exist")
	v.Check(util.MaxLen(item.Unit, 20), "unit", "Maximum length is 50")
	v.Check()
}
