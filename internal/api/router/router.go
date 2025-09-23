package router

import (
	"github.com/K1la/delayed-notifier/internal/api/handlers"
	"github.com/wb-go/wbf/ginext"
)

func New(handler *handlers.Handler) *ginext.Engine {
	e := ginext.New()
	e.Use(ginext.Recovery(), ginext.Logger())

	e.Static("/static", "./web")
	e.StaticFile("/", "./web/index.html")
	api := e.Group("/api/notify")
	{
		api.POST("/", handler.Create)
		api.GET("/", handler.GetAll)
		api.GET("/:id", handler.GetByID)
		api.DELETE("/:id", handler.Delete)
	}

	return e
}
