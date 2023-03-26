package main

import (
	"demo/logger"
	"demo/router"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.D.Debug("init")
	engine := gin.Default()
	router.RegisterAll(engine)
	engine.Run(":8060")
}
