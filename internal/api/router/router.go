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
		api.POST("/", handler.CreateNotification)
		api.GET("/", handler.GetAllNotifications)
		// TODO: доделать в ручках методы
		//api.GET("/:id", handler.GetNotificationStatusByID)
		//api.DELETE("/:id", handler.Delete)
	}

	return e
}
