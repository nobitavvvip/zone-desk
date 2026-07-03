package util

import "github.com/gin-gonic/gin"

func OK(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"success": true,
		"message": "ok",
		"data":    data,
	})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"message": message,
		"data":    nil,
	})
}
