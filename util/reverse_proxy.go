package util

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func CheckReverseProxy(c *gin.Context) {
	if !viper.GetBool("gin.reverse_proxy") {
		c.Request.Header.Del("Forwarded")
		c.Request.Header.Del("X-Forwarded-Proto")
		c.Request.Header.Del("X-Forwarded-Host")
		c.Request.Header.Del("X-Forwarded-For")
	}
}
