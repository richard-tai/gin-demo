package middleware

import (
	"demo/util/metrics"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover() func(c *gin.Context) {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		msg := "panic"
		if err != nil {
			if _, ok := err.(error); ok {
				msg = err.(error).Error()
			}
		}
		resp := map[string]interface{}{
			"error_message": msg,
		}
		c.JSON(http.StatusInternalServerError, resp)
	})
}

func CountApi(c *gin.Context) {
	metrics.ApiTotal.WithLabelValues(c.Request.URL.Path).Inc()
	c.Next()
}

func CheckPermisstion(c *gin.Context) {
	pass := false
	if pass == false {
		resp := map[string]interface{}{}
		c.JSON(http.StatusBadRequest, resp)
		c.Abort()
		return
	}
	c.Next()
}
