package main

import (
	router "lbe/api"
)

// @title           LBE API
// @version         1.0
// @description     Endpoints for authentication, login and register
// @description
// @description <details open>
// @description   <summary><a href="javascript:void(0)" style="cursor: pointer !important;">📋 Message Codes</a></summary>
// @description
// @description | Code   | Description                   |
// @description | ------ | ------------------------------|
// @description | 1000   | successful                    |
// @description | 1001   | unsuccessful                  |
// @description | 1002   | found                         |
// @description | 1003   | not found                     |
// @description | 4000   | internal error                |
// @description | 4001   | invalid request body          |
// @description | 4002   | invalid authentication token  |
// @description | 4003   | missing authentication token  |
// @description | 4004   | invalid signature             |
// @description | 4005   | missing signature             |
// @description | 4006   | invalid appid                 |
// @description | 4007   | missing appid                 |
// @description | 4008   | invalid query parameters      |
// @description
// @description </details>
// @host            localhost:18080
// @BasePath        /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
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
