package main

import (
	"CoinMarket/internal/config"
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/infrastructure/db/dsn"
	"CoinMarket/internal/infrastructure/repository"
	"gorm.io/gorm"
	"log"
)

func main() {
	log.Println("Migration started...")

	var err error

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := repository.CreateDB(dsn.FromCfg(cfg))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = RunMigrations(db)
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	log.Println("Migration completed successfully!")
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.InventoryItem{},
		&models.Transaction{},
		&models.Item{},
	)
}
