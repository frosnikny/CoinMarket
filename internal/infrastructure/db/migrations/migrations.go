package migrations

import (
	"CoinMarket/internal/domain/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.InventoryItem{},
		&models.Transaction{},
		&models.Item{},
	)
}
