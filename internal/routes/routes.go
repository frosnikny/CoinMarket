package routes

import (
	"CoinMarket/internal/app"
	"CoinMarket/internal/handlers"
	"CoinMarket/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes(r *gin.Engine, app *app.Application) {
	authHandler := handlers.NewAuthHandler(app.AuthService)
	walletHandler := handlers.NewWalletHandler(app.WalletService)

	api := r.Group("/api")
	{
		api.POST("/auth", authHandler.Auth)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(app.AuthService))

		protected.GET("/info", walletHandler.GetUserInfo)

		protected.POST("/sendCoin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Coin sent!"})
		})

		protected.GET("/buy/:item", func(c *gin.Context) {
			item := c.Param("item")
			c.JSON(http.StatusOK, gin.H{"message": "Item purchased!", "item": item})
		})
	}
}
