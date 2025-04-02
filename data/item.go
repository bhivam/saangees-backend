package data

import (
	"time"
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
