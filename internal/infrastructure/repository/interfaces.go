package repository

import (
	models2 "CoinMarket/internal/domain/models"
)

type WalletRepositoryInterface interface {
	GetBalance(username string) (int, error)
	UpdateBalance(username string, amount int) error
	GetInventory(username string) ([]models2.InventoryItem, error)
	GetTransactions(username string) ([]models2.Transaction, error)
	CreateTransaction(fromUsername, toUsername string, amount int) error
	UpdateInventory(username string, itemType string, quantity int) error
	AddToInventory(username string, itemType string, quantity int) error

	ExecuteDBTransaction(txFunc func(repo WalletRepositoryInterface) error) error
}

type UserRepositoryInterface interface {
	CreateUser(user *models2.User) error
	GetUserByUsername(username string) (*models2.User, error)
}

type ItemRepositoryInterface interface {
	GetItemByName(name string) (*models2.Item, error)
	SeedItems() error
}

type MainRepositoryInterface interface {
}
