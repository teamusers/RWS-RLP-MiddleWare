package main

import (
	router "rlp-member-service/api"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//topic.StartSubscription()

	router.Init()
}
