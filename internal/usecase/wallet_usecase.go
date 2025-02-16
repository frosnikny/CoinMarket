package usecase

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/infrastructure/repository"
	"errors"
)

type WalletService struct {
	WalletRepo repository.WalletRepositoryInterface
	UserRepo   repository.UserRepositoryInterface
	ItemRepo   repository.ItemRepositoryInterface
}

var _ WalletServiceInterface = (*WalletService)(nil)

func NewWalletService(walletRepo repository.WalletRepositoryInterface, userRepo repository.UserRepositoryInterface,
	itemRepo repository.ItemRepositoryInterface) *WalletService {
	return &WalletService{
		WalletRepo: walletRepo,
		UserRepo:   userRepo,
		ItemRepo:   itemRepo,
	}
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
	coins, err := s.WalletRepo.GetBalance(username)
	if err != nil {
		return nil, err
	}

	inventory, err := s.WalletRepo.GetInventory(username)
	if err != nil {
		return nil, err
	}

	transactions, err := s.WalletRepo.GetTransactions(username)
	if err != nil {
		return nil, err
	}

	// make для создания пустого массива
	received := make([]models.Transaction, 0)
	sent := make([]models.Transaction, 0)
	for _, transaction := range transactions {
		if transaction.ToUser == username {
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

func (s *WalletService) SendCoins(fromUserID, toUser string, amount int) error {
	if fromUserID == toUser {
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

	return s.WalletRepo.ExecuteDBTransaction(func(repo repository.WalletRepositoryInterface) error {
		if err := repo.UpdateBalance(fromUserID, -amount); err != nil {
			return err
		}

		if _, err := s.UserRepo.GetUserByUsername(toUser); err != nil {
			return err
		}

		if err := repo.UpdateBalance(toUser, amount); err != nil {
			return err
		}

		if err := repo.CreateTransaction(fromUserID, toUser, amount); err != nil {
			return err
		}

		return nil
	})
}

var (
	ErrItemNotFound   = errors.New("item not found")
	ErrNotEnoughCoins = errors.New("not enough coins to buy this item")
)

func (s *WalletService) BuyItem(user string, itemName string) error {
	item, err := s.ItemRepo.GetItemByName(itemName)
	if err != nil {
		return ErrItemNotFound
	}

	balance, err := s.WalletRepo.GetBalance(user)
	if err != nil {
		return err
	}
	if balance < item.Price {
		return ErrNotEnoughCoins
	}

	return s.WalletRepo.ExecuteDBTransaction(func(walletRepo repository.WalletRepositoryInterface) error {
		if err := walletRepo.UpdateBalance(user, -item.Price); err != nil {
			return err
		}

		inventory, err := walletRepo.GetInventory(user)
		if err != nil {
			return err
		}

		var found bool
		for _, invItem := range inventory {
			if invItem.ItemType == item.Name {
				found = true
				break
			}
		}

		if found {
			return walletRepo.UpdateInventory(user, item.Name, 1)
		} else {
			return walletRepo.AddToInventory(user, item.Name, 1)
		}
	})
}

func (s *WalletService) GetUserByUsername(username string) (*models.User, error) {
	return s.UserRepo.GetUserByUsername(username)
}
