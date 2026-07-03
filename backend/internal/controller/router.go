package controller

import (
	"zonedesk/internal/config"
	"zonedesk/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, shortcutSvc *service.ShortcutService) *gin.Engine {
	r := gin.Default()

	r.GET("/api/health", HealthHandler)
	r.GET("/api/config", ConfigHandler(cfg))

	files := r.Group("/api/files")
	{
		files.GET("/list", FileListHandler(cfg))
		files.GET("/stat", FileStatHandler(cfg))
		files.GET("/read", FileReadHandler(cfg))
		files.GET("/download", FileDownloadHandler(cfg))
		files.POST("/upload", FileUploadHandler(cfg))
		files.POST("/mkdir", FileMkdirHandler(cfg))
		files.POST("/rename", FileRenameHandler(cfg))
		files.POST("/delete", FileDeleteHandler(cfg))
		files.POST("/copy", FileCopyHandler(cfg))
		files.POST("/move", FileMoveHandler(cfg))
	}

	shortcuts := r.Group("/api/shortcuts")
	{
		shortcuts.GET("", ShortcutListHandler(shortcutSvc))
		shortcuts.POST("", ShortcutCreateHandler(shortcutSvc))
		shortcuts.DELETE("/:id", ShortcutDeleteHandler(shortcutSvc))
	}

	r.Static("/assets", cfg.Storage.WebDir+"/assets")
	r.NoRoute(func(c *gin.Context) {
		c.File(cfg.Storage.WebDir + "/index.html")
	})

	return r
}
