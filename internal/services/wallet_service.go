package services

import (
	"CoinMarket/internal/models"
	"CoinMarket/internal/repository"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepo *repository.WalletRepository
	UserRepo   *repository.UserRepository
}

func NewWalletService(walletRepo *repository.WalletRepository, userRepo *repository.UserRepository) *WalletService {
	return &WalletService{WalletRepo: walletRepo, UserRepo: userRepo}
}

type InfoResponse struct {
	Coins       int                    `json:"coins"`
	Inventory   []models.InventoryItem `json:"inventory"`
	CoinHistory struct {
		Received []models.Transaction `json:"received"`
		Sent     []models.Transaction `json:"sent"`
	} `json:"coinHistory"`
}

func (s *WalletService) GetUserInfo(userID uuid.UUID) (*InfoResponse, error) {
	coins, err := s.WalletRepo.GetBalance(userID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.WalletRepo.GetInventory(userID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.WalletRepo.GetTransactions(userID)
	if err != nil {
		return nil, err
	}

	// make для создания пустого массива
	received := make([]models.Transaction, 0)
	sent := make([]models.Transaction, 0)
	for _, transaction := range transactions {
		if transaction.ToUserID == userID {
			received = append(received, transaction)
		} else {
			sent = append(sent, transaction)
		}
	}

	response := &InfoResponse{
		Coins:     coins,
		Inventory: inventory,
	}
	response.CoinHistory.Received = received
	response.CoinHistory.Sent = sent

	return response, nil
}

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrNegativeAmount    = errors.New("amount must be positive")
	ErrSelfTransfer      = errors.New("cannot send coins to yourself")
)

func (s *WalletService) SendCoins(fromUserID, toUserID uuid.UUID, amount int) error {
	if fromUserID == toUserID {
		return ErrSelfTransfer
	}

	if amount <= 0 {
		return ErrNegativeAmount
	}

	balance, err := s.WalletRepo.GetBalance(fromUserID)
	if err != nil {
		return err
	}

	if balance < amount {
		return ErrInsufficientFunds
	}

	return s.WalletRepo.DB.Transaction(func(tx *gorm.DB) error {
		repo := repository.NewWalletRepository(tx)

		if err := repo.UpdateBalance(fromUserID, -amount); err != nil {
			return err
		}

		if err := repo.UpdateBalance(toUserID, amount); err != nil {
			return err
		}

		if err := repo.CreateTransaction(fromUserID, toUserID, amount); err != nil {
			return err
		}

		return nil
	})
}
