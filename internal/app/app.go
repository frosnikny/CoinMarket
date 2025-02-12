package app

import (
	"CoinMarket/internal/config"
	"CoinMarket/internal/dsn"
	"CoinMarket/internal/repository"
)

type Application struct {
	Repo   *repository.Repository
	Config *config.Config
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

	return a
}
