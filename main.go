package main

import (
	"demo/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.D.Debug("init")
	engine := gin.Default()
	engine.Run(":8060")
}
