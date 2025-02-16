package main

import (
	"CoinMarket/internal/app"
	"CoinMarket/internal/infrastructure/server"
	"log"
)

func main() {
	log.Println("Application start up!")
	a := app.New(nil)
	log.Println("Application created")
	server.StartServer(a)
	log.Println("Application terminated!")
}
