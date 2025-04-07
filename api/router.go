package router

import (
	"fmt"
	"net/http"

	general "rlp-middleware/api/http"
	"rlp-middleware/api/interceptor"
	"rlp-middleware/api/ws"
	"rlp-middleware/config"

	"github.com/gin-gonic/gin"
)

// Option type and global slice for router modifications.
type Option func(*gin.RouterGroup)

var options = []Option{}

var endpointList []map[string]string

func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init() *gin.Engine {
	// Include additional routers
	Include(general.Routers)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/index", helloHandler)

	wsGroup := r.Group("/ws", interceptor.WSInterceptor())
	wsGroup.GET("chat", ws.Chat)

	apiGroup := r.Group("/api", interceptor.HttpInterceptor())
	//apiGroup := r.Group("/api", interceptor.APITokenInterceptor())
	for _, opt := range options {
		opt(apiGroup)
	}

	// Capture routes but exclude the HandlerFunc to avoid JSON marshalling errors.
	routes := r.Routes()
	for _, route := range routes {
		endpointList = append(endpointList, map[string]string{
			"method": route.Method,
			"path":   route.Path,
		})
	}

	r.Run(fmt.Sprintf(":%d", config.GetConfig().Http.Port))
	return r
}

// helloHandler returns the filtered list of endpoints in JSON.
func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, endpointList)
}
