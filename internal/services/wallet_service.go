package services

import (
	"CoinMarket/internal/models"
	"CoinMarket/internal/repository"
)

type WalletService struct {
	repo *repository.WalletRepository
}

func NewWalletService(repo *repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

type InfoResponse struct {
	Coins       int                    `json:"coins"`
	Inventory   []models.InventoryItem `json:"inventory"`
	CoinHistory struct {
		Received []models.Transaction `json:"received"`
		Sent     []models.Transaction `json:"sent"`
	} `json:"coinHistory"`
}

func (s *WalletService) GetUserInfo(username string) (*InfoResponse, error) {
	coins, err := s.repo.GetBalance(username)
	if err != nil {
		return nil, err
	}

	inventory, err := s.repo.GetInventory(username)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.GetTransactions(username)
	if err != nil {
		return nil, err
	}

	var received, sent []models.Transaction
	for _, tx := range transactions {
		if tx.ToUser == username {
			received = append(received, tx)
		} else {
			sent = append(sent, tx)
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
