package demo_v1_handler

import (
	"demo/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ConfsGet(c *gin.Context) {
	logger.D.Debug("URL [%+v]", *(c.Request.URL))
	logger.D.Debug("Host [%+v]", c.Request.Host)
	logger.D.Debug("Method [%+v]", c.Request.Method)
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}
