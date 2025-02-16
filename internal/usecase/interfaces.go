package usecase

import (
	"CoinMarket/internal/domain/models"
)

type WalletServiceInterface interface {
	GetUserInfo(username string) (*InfoResponse, error)
	SendCoins(fromUser, toUser string, amount int) error
	BuyItem(username string, itemName string) error
	GetUserByUsername(username string) (*models.User, error)
}
