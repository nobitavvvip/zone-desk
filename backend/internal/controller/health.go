package controller

import (
	"zonedesk/pkg/util"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	util.OK(c, gin.H{
		"status":  "running",
		"version": "0.1.0",
	})
}
