package server

import (
	"CoinMarket/internal/app"
	"CoinMarket/internal/routes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartServer(a *app.Application) {
	log.Println("Server start up")

	r := gin.Default()

	// r.Use(ErrorHandler())

	routes.SetupRoutes(r, a)

	s := &http.Server{
		Addr:    a.Config.ServerAddress,
		Handler: r,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}

	log.Println("Server down")
}
