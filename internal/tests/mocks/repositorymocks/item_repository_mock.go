package repositorymocks

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/infrastructure/repository"
	"github.com/stretchr/testify/mock"
)

type ItemRepositoryMock struct {
	mock.Mock
}

var _ repository.ItemRepositoryInterface = (*ItemRepositoryMock)(nil)

func (m *ItemRepositoryMock) GetItemByName(name string) (*models.Item, error) {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Item), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ItemRepositoryMock) SeedItems() error {
	args := m.Called()
	return args.Error(0)
}
