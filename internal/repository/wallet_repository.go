package repository

import (
	"CoinMarket/internal/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetBalance(username string) (int, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user.Coins, err
}

func (r *WalletRepository) GetTransactions(username string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("from_user = ? OR to_user = ?", username, username).Find(&transactions).Error
	return transactions, err
}

func (r *WalletRepository) GetInventory(username string) ([]models.InventoryItem, error) {
	var inventory []models.InventoryItem
	err := r.db.Where("username = ?", username).Find(&inventory).Error
	return inventory, err
}
