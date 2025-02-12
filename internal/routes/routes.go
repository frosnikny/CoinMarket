package routes

import (
	"CoinMarket/internal/app"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, app *app.Application) {
	_ = r.Group("/api")
	{
	}
}
