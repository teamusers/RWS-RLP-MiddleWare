package router

import (
	"fmt"

	general "lbe/api/http"
	"lbe/config"

	"github.com/gin-gonic/gin"

	// <-- this import makes sure your docs/swagger.json is registered
	_ "lbe/api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Option func(*gin.RouterGroup)

var options = []Option{}
var endpointList []map[string]string

func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init() *gin.Engine {
	Include(general.Routers)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// mount your API routes
	apiGroup := r.Group("/api")
	for _, opt := range options {
		opt(apiGroup)
	}

	// capture endpoints if you need them
	for _, route := range r.Routes() {
		endpointList = append(endpointList, map[string]string{
			"method": route.Method,
			"path":   route.Path,
		})
	}

	// in router.go, after you set up your API routes…

	// serve everything in ./api/docs at the URL path /docs
	r.Static("/docs", "./api/docs")

	// wire up the swagger UI, telling it to fetch /docs/swagger.json
	url := ginSwagger.URL("/docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// start server
	r.Run(fmt.Sprintf(":%d", config.GetConfig().Http.Port))
	return r
}
