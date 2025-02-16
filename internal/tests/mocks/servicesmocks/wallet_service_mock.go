package servicesmocks

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/usecase"
	"github.com/stretchr/testify/mock"
)

type WalletServiceMock struct {
	mock.Mock
}

var _ usecase.WalletServiceInterface = (*WalletServiceMock)(nil)

func (m *WalletServiceMock) GetUserInfo(username string) (*usecase.InfoResponse, error) {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*usecase.InfoResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *WalletServiceMock) SendCoins(fromUsername, toUsername string, amount int) error {
	args := m.Called(fromUsername, toUsername, amount)
	return args.Error(0)
}

func (m *WalletServiceMock) BuyItem(username string, itemName string) error {
	args := m.Called(username, itemName)
	return args.Error(0)
}

func (m *WalletServiceMock) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}
