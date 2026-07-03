package controller

import (
	"zonedesk/internal/config"
	"zonedesk/pkg/util"

	"github.com/gin-gonic/gin"
)

func ConfigHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		util.OK(c, gin.H{
			"fileManager": gin.H{
				"defaultPath": cfg.FileManager.DefaultPath,
			},
			"ui": gin.H{
				"theme": cfg.UI.Theme,
			},
		})
	}
}
