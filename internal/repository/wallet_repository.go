package repository

import (
	"CoinMarket/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (r *WalletRepository) GetTransactions(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Find(&transactions).Error
	return transactions, err
}

func (r *WalletRepository) GetInventory(userID uuid.UUID) ([]models.InventoryItem, error) {
	var inventory []models.InventoryItem
	err := r.DB.Where("user_id = ?", userID).Find(&inventory).Error
	return inventory, err
}

func (r *WalletRepository) UpdateBalance(userID uuid.UUID, amount int) error {
	return r.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("coins", gorm.Expr("coins + ?", amount)).Error
}

func (r *WalletRepository) CreateTransaction(fromUserID, toUserID uuid.UUID, amount int) error {
	tx := models.Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
	}
	return r.DB.Create(&tx).Error
}

func (r *WalletRepository) GetBalance(userID uuid.UUID) (int, error) {
	var user models.User
	err := r.DB.Select("coins").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.Coins, nil
}
