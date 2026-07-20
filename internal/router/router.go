package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	todohandler "github.com/arrebole/vibe-coding-starter/internal/modules/todo/handler"
	"github.com/arrebole/vibe-coding-starter/internal/response"
)

type Dependencies struct {
	DB          *gorm.DB
	TodoHandler *todohandler.Handler
}

func New(deps Dependencies) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})

	api := engine.Group("/api/v1")
	deps.TodoHandler.Register(api)

	engine.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "接口不存在")
	})

	return engine
}
