package app

import (
	"CoinMarket/internal/config"
	"CoinMarket/internal/dsn"
	"CoinMarket/internal/repository"
	"CoinMarket/internal/services"
)

type Application struct {
	Repo          *repository.Repository
	UserRepo      *repository.UserRepository
	WalletRepo    *repository.WalletRepository
	AuthService   *services.AuthService
	WalletService *services.WalletService
	Config        *config.Config
}

func New() *Application {
	var err error

	a := &Application{}
	a.Config, err = config.New()
	if err != nil {
		panic(err)
	}

	a.Repo, err = repository.New(dsn.FromCfg(a.Config))
	if err != nil {
		panic(err)
	}

	a.UserRepo = repository.NewUserRepository(a.Repo.DB)
	a.WalletRepo = repository.NewWalletRepository(a.Repo.DB)

	a.AuthService = services.NewAuthService(a.UserRepo, a.Config.JwtKey)
	a.WalletService = services.NewWalletService(a.WalletRepo)

	return a
}
