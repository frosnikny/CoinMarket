package repository_mocks

import (
	"CoinMarket/internal/infrastructure/repository"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

var _ repository.MainRepositoryInterface = (*RepositoryMock)(nil)
