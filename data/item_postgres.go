package data

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type PostgresItemStore struct {
	DB *gorm.DB
}

func NewPostgresItemStore(db *gorm.DB) *PostgresItemStore {
	if err := db.AutoMigrate(&Item{}); err != nil {
		panic("failed to migrate item schema: " + err.Error())
	}
	if err := db.AutoMigrate(&ModifierCategory{}); err != nil {
		panic("failed to migrate modifier category schema: " + err.Error())
	}
	if err := db.AutoMigrate(&ModifierOption{}); err != nil {
		panic("failed to migrate modifier option schema: " + err.Error())
	}
	return &PostgresItemStore{DB: db}
}

func (store *PostgresItemStore) CreateItem(item *Item) (*Item, error) {
	if err := store.DB.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (store *PostgresItemStore) GetItem(id uint) (*Item, error) {
	var item Item
	if err := store.DB.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}
	return &item, nil
}

// TODO check if this works as expected
// TODO add everywhere else if it does
func (store *PostgresItemStore) UpdateItem(item *Item) error {
	if err := store.DB.Updates(item).Error; err != nil {
		return err
	}
	return nil
}

func (store *PostgresItemStore) DeleteItem(id uint) error {
	if err := store.DB.Delete(&Item{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (store *PostgresItemStore) ComingWeekItems() ([]*Item, error) {
	var items []*Item
	oneWeekFromNow := time.Now().AddDate(0, 0, 7)
	if err := store.DB.Where("date BETWEEN ? AND ?", time.Now(), oneWeekFromNow).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
