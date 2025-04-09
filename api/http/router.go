package http

import (
	v1 "rlp-member-service/api/http/controllers/v1"
	"rlp-member-service/api/interceptor"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {

	v1Group := e.Group("/v1")

	v1Group.GET("/auth", v1.AuthHandler)

	// Create a sub-group for "/users" with the HttpInterceptor middleware applied.
	usersGroup := v1Group.Group("/user", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		usersGroup.GET("", v1.GetUser)
		usersGroup.POST("", v1.CreateUser)
		//usersGroup.PUT("/:id", v1.UpdateUser)
		//usersGroup.DELETE("/:id", v1.DeleteUser)
	}

	loginGroup := v1Group.Group("/login/user", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		loginGroup.GET("", v1.Login)
		//usersGroup.PUT("/:id", v1.UpdateUser)
		//usersGroup.DELETE("/:id", v1.DeleteUser)
	}

}
