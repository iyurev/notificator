package handlers

import "github.com/gin-gonic/gin"

func ContentTypeIsJSON(c *gin.Context) bool {
	if c.GetHeader("Content-Type") == "application/json" {
		return true
	}
	return false
}
