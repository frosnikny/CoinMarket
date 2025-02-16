package repository

import (
	"CoinMarket/internal/domain/models"
	"errors"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

var _ ItemRepositoryInterface = (*ItemRepository)(nil)

func NewItemRepository(db *gorm.DB) (*ItemRepository, error) {
	itemRepo := ItemRepository{db: db}
	err := itemRepo.SeedItems()
	if err != nil {
		return nil, err
	}
	return &itemRepo, nil
}

func (r *ItemRepository) GetItemByName(name string) (*models.Item, error) {
	var item models.Item
	err := r.db.Where("name = ?", name).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("item not found")
	}
	return &item, err
}

func (r *ItemRepository) SeedItems() error {
	items := []models.Item{
		{Name: "t-shirt", Price: 80},
		{Name: "cup", Price: 20},
		{Name: "book", Price: 50},
		{Name: "pen", Price: 10},
		{Name: "powerbank", Price: 200},
		{Name: "hoody", Price: 300},
		{Name: "umbrella", Price: 200},
		{Name: "socks", Price: 10},
		{Name: "wallet", Price: 50},
		{Name: "pink-hoody", Price: 500},
	}

	for _, item := range items {
		if _, err := r.GetItemByName(item.Name); err != nil {
			r.db.Create(&item)
		}
	}

	return nil
}
