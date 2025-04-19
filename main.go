package main

import (
	router "lbe/api"
)

// @title           LBE API
// @version         1.0
// @description     Endpoints for authentication, login and register
// @host            localhost:18080
// @BasePath        /api/v1

// @securityDefinitions.apikey  ApiKeyAuth  // arbitrary name
// @in                         header
// @name                       Authorization
// @description                Type "Bearer <your-jwt>" to authorize
func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//topic.StartSubscription()
	//gin.SetMode(gin.ReleaseMode)
	router.Init()
}
