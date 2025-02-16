package repositorymocks

import (
	"CoinMarket/internal/domain/models"
	"CoinMarket/internal/infrastructure/repository"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

var _ repository.UserRepositoryInterface = (*UserRepositoryMock)(nil)

func (m *UserRepositoryMock) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}
