package dsn

import (
	"CoinMarket/internal/config"
	"fmt"
)

func FromCfg(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUsername, cfg.PostgresPassword, cfg.PostgresDatabase)
}
