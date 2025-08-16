package router

import (
	"fmt"
	"net/url"
	"time"

	"github.com/amahdian/golang-gin-boilerplate/svc/auth"

	"github.com/amahdian/golang-gin-boilerplate/docs"
	"github.com/amahdian/golang-gin-boilerplate/global/env"
	"github.com/amahdian/golang-gin-boilerplate/pkg/logger"
	"github.com/amahdian/golang-gin-boilerplate/server/binding"
	"github.com/amahdian/golang-gin-boilerplate/server/middleware"
	"github.com/amahdian/golang-gin-boilerplate/storage"
	"github.com/amahdian/golang-gin-boilerplate/svc"
	"github.com/amahdian/golang-gin-boilerplate/version"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	*gin.Engine

	storage storage.Storage
	svc     svc.Svc

	configs *env.Envs

	authenticator auth.Authenticator

	publicGroup *gin.RouterGroup
	authGroup   *gin.RouterGroup
	apiGroup    *gin.RouterGroup
}

func NewRouter(
	storage storage.Storage,
	svc svc.Svc,
	configs *env.Envs,
	authenticator auth.Authenticator) *Router {
	gin.SetMode(configs.Server.GinMode)
	router := &Router{
		Engine:        gin.New(),
		storage:       storage,
		authenticator: authenticator,
		svc:           svc,
		configs:       configs,
	}
	router.Use(otelgin.Middleware(version.AppName))
	router.Use(
		middleware.WithLogger(),
		middleware.WithRecovery(),
	)
	pprof.Register(router.Engine)
	router.setupBindings()
	router.setupCors()
	router.setupSwagger()
	router.setupRoutes()

	return router
}

func (r *Router) setupBindings() {
	binding.Init()
}

func (r *Router) setupCors() {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "HEAD", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))
}

func (r *Router) setupSwagger() {
	scheme := "http"
	host := fmt.Sprintf("localhost:%s", r.configs.Server.HttpPort)
	basePath := "/"

	if r.configs.Server.SwaggerHostAddr != "" {
		uri, err := url.ParseRequestURI(r.configs.Server.SwaggerHostAddr)
		if err != nil {
			logger.Errorf("failed to parse swagger host address: %v", err)
		} else {
			scheme = uri.Scheme
			host = uri.Host
			basePath = uri.Path
		}
	}

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Schemes = []string{scheme}
}
