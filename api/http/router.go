package http

import (
	v1 "rlp-middleware/api/http/controller/v1"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {

	v1Group := e.Group("/v1")

	v1Group.GET("/users", v1.GetUsers)
	v1Group.GET("/users/:id", v1.GetUser)
	v1Group.POST("/users", v1.CreateUser)
	v1Group.PUT("/users/:id", v1.UpdateUser)
	v1Group.DELETE("/users/:id", v1.DeleteUser)

}
