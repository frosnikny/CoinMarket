package main

import (
	"CoinMarket/internal/app"
	"CoinMarket/internal/server"
	"log"
)

func main() {
	log.Println("Application start up!")
	a := app.New()
	log.Println("Application created")
	server.StartServer(a)
	log.Println("Application terminated!")
}
