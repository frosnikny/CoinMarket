package repository

import (
	"CoinMarket/internal/domain/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

var _ WalletRepositoryInterface = (*WalletRepository)(nil)

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (r *WalletRepository) GetTransactions(username string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("from_user = ? OR to_user = ?", username, username).Find(&transactions).Error
	return transactions, err
}

func (r *WalletRepository) GetInventory(username string) ([]models.InventoryItem, error) {
	var inventory []models.InventoryItem
	err := r.DB.Where("username = ?", username).Find(&inventory).Error
	return inventory, err
}

func (r *WalletRepository) UpdateBalance(username string, amount int) error {
	return r.DB.Model(&models.User{}).Where("username = ?", username).
		Update("coins", gorm.Expr("coins + ?", amount)).Error
}

func (r *WalletRepository) CreateTransaction(fromUsername, toUsername string, amount int) error {
	tx := models.Transaction{
		FromUser: fromUsername,
		ToUser:   toUsername,
		Amount:   amount,
	}
	return r.DB.Create(&tx).Error
}

func (r *WalletRepository) GetBalance(username string) (int, error) {
	var user models.User
	err := r.DB.Select("coins").Where("username = ?", username).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.Coins, nil
}

func (r *WalletRepository) UpdateInventory(username string, itemType string, quantity int) error {
	return r.DB.Model(&models.InventoryItem{}).
		Where("username = ? AND item_type = ?", username, itemType).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
}

func (r *WalletRepository) AddToInventory(username string, itemType string, quantity int) error {
	newItem := models.InventoryItem{Username: username, ItemType: itemType, Quantity: quantity}
	return r.DB.Create(&newItem).Error
}

func (r *WalletRepository) ExecuteDBTransaction(txFunc func(repo WalletRepositoryInterface) error) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		repo := NewWalletRepository(tx)
		return txFunc(repo)
	})
}
