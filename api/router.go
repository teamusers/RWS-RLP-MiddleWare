package router

import (
	"fmt"
	"net/http"

	general "github.com/stonksdex/externalapi/api/http"
	"github.com/stonksdex/externalapi/api/interceptor"
	"github.com/stonksdex/externalapi/api/ws"
	"github.com/stonksdex/externalapi/config"

	"github.com/gin-gonic/gin"
)

type Option func(*gin.RouterGroup)

var options = []Option{}

func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init() *gin.Engine {
	Include(general.Routers)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/index", helloHandler) //Default welcome api

	wsGroup := r.Group("/ws", interceptor.WSInterceptor())
	wsGroup.GET("chat", ws.Chat)

	apiGroup := r.Group("/spwapi", interceptor.HttpInterceptor()) // total interceptor stack
	for _, opt := range options {
		opt(apiGroup)
	}
	r.Run(fmt.Sprintf(":%d", config.GetConfig().Http.Port))
	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello Stonks",
	})
}
