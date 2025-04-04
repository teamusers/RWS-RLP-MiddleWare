package main

import (
	"github.com/joho/godotenv"
	router "github.com/stonksdex/externalapi/api"
	"github.com/stonksdex/externalapi/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//topic.StartSubscription()

	router.Init()
}
