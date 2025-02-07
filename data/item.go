package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type SizeOption struct {
	Name          string  `json:"name"`
	PriceModifier float64 `json:"price_modifier"`
}

type SpiceOption struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SizeOptions []SizeOption

func (o SizeOptions) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *SizeOptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &o)
}

type SpiceOptions []SpiceOption

func (o SpiceOptions) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *SpiceOptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &o)
}

type Item struct {
	ID           int64        `json:"id"            gorm:"primaryKey"`
	Name         string       `json:"name"          gorm:"type:text;not null"`
	Description  string       `json:"description"   gorm:"type:text;not null"`
	BasePrice    float64      `json:"base_price"    gorm:"type:decimal(10,2);not null"`
	Date         time.Time    `json:"date"          gorm:"type:date;not null"`
	SizeOptions  SizeOptions  `json:"size_options"  gorm:"type:jsonb;not null;default:'[]'"`
	SpiceOptions SpiceOptions `json:"spice_options" gorm:"type:jsonb;not null;default:'[]'"`
}

type ItemStore interface {
	CreateItem(item *Item) (*Item, error)
	GetItem(id int64) (*Item, error)
	UpdateItem(item *Item) error
	DeleteItem(id int64) error
	ComingWeekItems() ([]*Item, error)
}
