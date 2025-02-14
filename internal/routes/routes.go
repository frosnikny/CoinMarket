package routes

import (
	"CoinMarket/internal/app"
	"CoinMarket/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, app *app.Application) {
	authHandler := handlers.NewAuthHandler(app.AuthService)

	api := r.Group("/api")
	{
		api.POST("/auth", authHandler.Auth)
	}
}
