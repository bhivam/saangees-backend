package data

import "time"

type Item struct {
	ID   int64     `json:"id"   gorm:"primaryKey"`
	Name string    `json:"name" gorm:"type:text;not null"`
	Date time.Time `json:"date" gorm:"type:date;not null"`
}

type ItemStore interface {
	CreateItem(item *Item) (*Item, error)
	GetItem(id int64) (*Item, error)
	UpdateItem(item *Item) error
	DeleteItem(id int64) error
	ComingWeekItems() ([]*Item, error)
}
