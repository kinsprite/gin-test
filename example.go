package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
)

func main() {
	engine := gin.New()
	engine.Use(apmgin.Middleware(engine))

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	engine.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
