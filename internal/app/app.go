package app

import (
	"CoinMarket/internal/config"
	"CoinMarket/internal/infrastructure/db/dsn"
	"CoinMarket/internal/infrastructure/repository"
	"CoinMarket/internal/usecase"
)

type Application struct {
	Repo       *repository.Repository
	UserRepo   *repository.UserRepository
	WalletRepo *repository.WalletRepository
	ItemRepo   *repository.ItemRepository

	AuthService   *usecase.AuthService
	WalletService *usecase.WalletService
	Config        *config.Config
}

func New(readyConfig *config.Config) *Application {
	var err error

	a := &Application{}
	if readyConfig != nil {
		a.Config = readyConfig
	} else {
		a.Config, err = config.New()
	}
	if err != nil {
		panic(err)
	}

	a.Repo, err = repository.New(dsn.FromCfg(a.Config))
	if err != nil {
		panic(err)
	}

	a.UserRepo = repository.NewUserRepository(a.Repo.DB)
	a.WalletRepo = repository.NewWalletRepository(a.Repo.DB)
	a.ItemRepo, err = repository.NewItemRepository(a.Repo.DB)
	if err != nil {
		panic(err)
	}

	a.AuthService = usecase.NewAuthService(a.UserRepo, a.Config.JwtKey)
	a.WalletService = usecase.NewWalletService(a.WalletRepo, a.UserRepo, a.ItemRepo)

	return a
}
