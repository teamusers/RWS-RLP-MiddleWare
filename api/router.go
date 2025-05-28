package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	general "lbe/api/http"
	"lbe/api/http/middleware"
	"lbe/config"
	"lbe/model"
	"lbe/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	var db *gorm.DB
	if os.Getenv("RUN_UNIT_TESTS") == "true" {
		log.Println("🧪  Test mode: skipping migrations")
	} else {
		// 0) run audit‐table migration
		db = system.GetDb()
		if err := model.MigrateAuditLog(db); err != nil {
			log.Fatalf("audit log migration: %v", err)
		}
		if err := model.MigrateRLPUserNumbering(db); err != nil {
			log.Fatalf("rlp user numbering migration: %v", err)
		}
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	httpClient := &http.Client{Timeout: 10 * time.Second}
	r.Use(middleware.HttpClientMiddleware(httpClient))

	// only wire AuditLogger if we have a real DB
	if db != nil {
		r.Use(middleware.AuditLogger(db))
	}

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		url,
		// ⟵ hides the Models section
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// redirect root to swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})
	// also catch bare /swagger
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})

	r.Run(fmt.Sprintf(":%d", config.GetConfig().Http.Port))
	return r
}
