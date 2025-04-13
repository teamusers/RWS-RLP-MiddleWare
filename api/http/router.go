package http

import (
	v1 "lbe/api/http/controllers/v1"
	user "lbe/api/http/controllers/v1/user"
	"lbe/api/interceptor"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {

	v1Group := e.Group("/v1")

	v1Group.GET("/auth", v1.AuthHandler)

	// Create a sub-group for "/users" with the HttpInterceptor middleware applied.
	usersGroup := v1Group.Group("/user/register", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		usersGroup.GET("/:email/:sign_up_type", user.GetUser)
		usersGroup.POST("", user.CreateUser)
		//usersGroup.PUT("/:id", v1.UpdateUser)
		//usersGroup.DELETE("/:id", v1.DeleteUser)
	}

	loginGroup := v1Group.Group("/user/login", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		loginGroup.GET("/:email", user.Login)
		//usersGroup.PUT("/:id", v1.UpdateUser)
		//usersGroup.DELETE("/:id", v1.DeleteUser)
	}

}
