package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iyurev/notificator/pkg/gitlab"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/gitlab", gitlab.WebHookHandler())
	r.Run()
}
