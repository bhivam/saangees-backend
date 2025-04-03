package data

import (
	"fmt"
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

func ValidateModifierOptions(v *util.Validator, mo []ModifierOption, prefix string) {

		for i, option := range mo {
			optionPrefix := fmt.Sprintf("%smodifier_options[%d].", prefix, i)

			v.Check(util.Nonempty(option.Name), optionPrefix+"name", "Must exist")
			v.Check(util.MaxLen(option.Name, 50), optionPrefix+"name", "Maximum length is 50")
      v.Check(util.NonNegativeFl(option.PriceModifier), optionPrefix+"price_modifier", "Must be nonnegative")
		}

}

func ValidateModifierCategories(v *util.Validator, mc []ModifierCategory) {
	for i, m := range mc {
		prefix := fmt.Sprintf("modifier_categories[%d].", i)

		v.Check(util.Nonempty(m.Name), prefix+"name", "Must exist")
		v.Check(util.MaxLen(m.Name, 50), prefix+"name", "Maximum length is 50")
		v.Check(m.Min >= 0, prefix+"min", "Must be non-negative")
		v.Check(m.Max >= m.Min, prefix+"max", "Must be greater than or equal to min")
		v.Check(
			len(m.ModifierOptions) > 0,
			prefix+"modifier_options",
			"Must have at least one option",
		)



		var optionNames []string
		for _, option := range m.ModifierOptions {
			optionNames = append(optionNames, option.Name)
		}
		v.Check(util.Unique(optionNames), prefix+"modifier_options", "Option names must be unique")
	}
}

func ValidateItem(v *util.Validator, item *Item) {
	v.Check(util.Nonempty(item.Name), "name", "Must exist")
	v.Check(util.MaxLen(item.Name, 50), "name", "Maximum length is 50")
	v.Check(util.Nonempty(item.Description), "description", "Must exist")
	v.Check(util.MaxLen(item.Description, 300), "description", "Maximum length is 300")
	v.Check(util.NonNegativeFl(item.BasePrice), "base_price", "Must be nonnegative")
	v.Check(item.Quantity > 0, "quantity", "Must be greater than 0")
	v.Check(util.Nonempty(item.Unit), "unit", "Must exist")
	v.Check(util.MaxLen(item.Unit, 20), "unit", "Maximum length is 20")

	if len(item.ModifierCategories) > 0 {
		ValidateModifierCategories(v, item.ModifierCategories)
	}

	var categoryNames []string
	for _, category := range item.ModifierCategories {
		categoryNames = append(categoryNames, category.Name)
	}
	v.Check(util.Unique(categoryNames), "modifier_categories", "Category names must be unique")
}
