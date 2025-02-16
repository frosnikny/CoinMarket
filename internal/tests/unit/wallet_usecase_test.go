package unit

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/infrastructure/repository"
	"CoinMarket/internal/tests/mocks/repository_mocks"
	"CoinMarket/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	mockWalletRepo := new(repository_mocks.WalletRepositoryMock)
	mockUserRepo := new(repository_mocks.UserRepositoryMock)
	mockItemRepo := new(repository_mocks.ItemRepositoryMock)

	service := usecase.NewWalletService(mockWalletRepo, mockUserRepo, mockItemRepo)

	username := "test_user"
	mockWalletRepo.On("GetBalance", username).Return(1000, nil)
	mockWalletRepo.On("GetInventory", username).Return([]models.InventoryItem{}, nil)
	mockWalletRepo.On("GetTransactions", username).Return([]models.Transaction{}, nil)

	info, err := service.GetUserInfo(username)
	assert.NoError(t, err)
	assert.Equal(t, 1000, info.Coins)
	assert.Empty(t, info.Inventory)
	assert.Empty(t, info.CoinHistory.Received)
	assert.Empty(t, info.CoinHistory.Sent)

	mockWalletRepo.AssertExpectations(t)
}

func TestSendCoins(t *testing.T) {
	mockWalletRepo := new(repository_mocks.WalletRepositoryMock)
	mockUserRepo := new(repository_mocks.UserRepositoryMock)
	mockItemRepo := new(repository_mocks.ItemRepositoryMock)

	service := usecase.NewWalletService(mockWalletRepo, mockUserRepo, mockItemRepo)

	fromUsername := "test_from_username"
	toUsername := "test_to_username"
	amount := 100

	mockUserRepo.On("GetUserByUsername", toUsername).Return(&models.User{Username: toUsername}, nil)

	mockWalletRepo.On("GetBalance", fromUsername).Return(200, nil)

	mockWalletRepo.On("ExecuteDBTransaction", mock.Anything).Run(func(args mock.Arguments) {
		txFunc := args.Get(0).(func(repo repository.WalletRepositoryInterface) error)
		_ = txFunc(mockWalletRepo)
	}).Return(nil)

	mockWalletRepo.On("UpdateBalance", fromUsername, -amount).Return(nil)
	mockWalletRepo.On("UpdateBalance", toUsername, amount).Return(nil)
	mockWalletRepo.On("CreateTransaction", fromUsername, toUsername, amount).Return(nil)

	err := service.SendCoins(fromUsername, toUsername, amount)

	assert.NoError(t, err)

	mockWalletRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBuyItem(t *testing.T) {
	mockWalletRepo := new(repository_mocks.WalletRepositoryMock)
	mockUserRepo := new(repository_mocks.UserRepositoryMock)
	mockItemRepo := new(repository_mocks.ItemRepositoryMock)

	service := usecase.NewWalletService(mockWalletRepo, mockUserRepo, mockItemRepo)

	username := "test_username"
	itemName := "t-shirt"

	mockItemRepo.On("GetItemByName", itemName).Return(&models.Item{Name: itemName, Price: 80}, nil)
	mockWalletRepo.On("GetBalance", username).Return(1000, nil)

	mockWalletRepo.On("ExecuteDBTransaction", mock.Anything).Run(func(args mock.Arguments) {
		txFunc := args.Get(0).(func(repo repository.WalletRepositoryInterface) error)
		_ = txFunc(mockWalletRepo)
	}).Return(nil)

	mockWalletRepo.On("UpdateBalance", username, -80).Return(nil)
	mockWalletRepo.On("GetInventory", username).Return([]models.InventoryItem{}, nil)
	mockWalletRepo.On("AddToInventory", username, itemName, 1).Return(nil)

	err := service.BuyItem(username, itemName)
	assert.NoError(t, err)

	mockWalletRepo.AssertExpectations(t)
	mockItemRepo.AssertExpectations(t)
}
