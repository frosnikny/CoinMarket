package main

import (
	"CoinMarket/internal/config"
	"CoinMarket/internal/dsn"
	"CoinMarket/internal/models"
	"CoinMarket/internal/repository"
	"gorm.io/gorm"
	"log"
)

func main() {
	log.Println("Миграция началась...")

	var err error

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := repository.CreateDB(dsn.FromCfg(cfg))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	log.Println("Миграция завершена успешно!")
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.InventoryItem{},
		&models.Transaction{},
	)
}
