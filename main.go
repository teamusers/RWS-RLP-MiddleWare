package main

import (
	router "rlp-middleware/api"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//topic.StartSubscription()

	router.Init()
}
