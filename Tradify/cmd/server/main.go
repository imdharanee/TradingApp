package main

import (
	"fmt"
	"log"

	"dex/internal/api"
	"dex/internal/engine"
	"dex/internal/orderbook"

	"github.com/gin-gonic/gin"
)

func main() {

	book := orderbook.NewOrderBook()

	matchingEngine := engine.NewMatchingEngine(book)

	go matchingEngine.Start()

	handler := api.NewHandler(book, matchingEngine)

	router := gin.Default()

	handler.SetupRoutes(router)

	fmt.Println("DEX server starting on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
