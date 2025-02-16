package routes

import (
	"CoinMarket/internal/app"
	"CoinMarket/internal/delivery/handlers"
	"CoinMarket/internal/infrastructure/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, app *app.Application) {
	authHandler := handlers.NewAuthHandler(app.AuthService)
	walletHandler := handlers.NewWalletHandler(app.WalletService)

	api := r.Group("/api")
	{
		api.POST("/auth", authHandler.Auth)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(app.AuthService, app.UserRepo))

		protected.GET("/info", walletHandler.GetUserInfo)
		protected.POST("/sendCoin", walletHandler.SendCoins)
		protected.GET("/buy/:item", walletHandler.BuyItem)
	}
}
