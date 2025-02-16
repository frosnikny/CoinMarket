package repository_mocks

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/infrastructure/repository"
	"github.com/stretchr/testify/mock"
)

type WalletRepositoryMock struct {
	mock.Mock
}

var _ repository.WalletRepositoryInterface = (*WalletRepositoryMock)(nil)

func (m *WalletRepositoryMock) GetBalance(username string) (int, error) {
	args := m.Called(username)
	return args.Int(0), args.Error(1)
}

func (m *WalletRepositoryMock) UpdateBalance(username string, amount int) error {
	args := m.Called(username, amount)
	return args.Error(0)
}

func (m *WalletRepositoryMock) GetInventory(username string) ([]models.InventoryItem, error) {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).([]models.InventoryItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *WalletRepositoryMock) GetTransactions(username string) ([]models.Transaction, error) {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *WalletRepositoryMock) CreateTransaction(fromUsername, toUsername string, amount int) error {
	args := m.Called(fromUsername, toUsername, amount)
	return args.Error(0)
}

func (m *WalletRepositoryMock) UpdateInventory(username string, itemType string, quantity int) error {
	args := m.Called(username, itemType, quantity)
	return args.Error(0)
}

func (m *WalletRepositoryMock) AddToInventory(username string, itemType string, quantity int) error {
	args := m.Called(username, itemType, quantity)
	return args.Error(0)
}

func (m *WalletRepositoryMock) ExecuteDBTransaction(txFunc func(repo repository.WalletRepositoryInterface) error) error {
	args := m.Called(txFunc)
	return args.Error(0)
}
