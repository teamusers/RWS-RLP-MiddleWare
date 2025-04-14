package router

import (
	"fmt"

	general "lbe/api/http"
	"lbe/config"

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
	apiGroup := r.Group("/api")
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
