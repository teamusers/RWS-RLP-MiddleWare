package main

import (
	router "lbe/api"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//topic.StartSubscription()
	//gin.SetMode(gin.ReleaseMode)
	router.Init()
}
