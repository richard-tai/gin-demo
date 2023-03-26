package router

import (
	"demo/handler/demo_v1_handler"
	"demo/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAll(e *gin.Engine) {
	e.Use(middleware.Recover())
	e.Use(middleware.CountApi)
	demo := e.Group("/demo")
	{
		v1 := demo.Group("/v1")
		{
			v1.GET("/confs", demo_v1_handler.ConfsGet)
		}
	}
	e.Static("/static", "./static")
}
