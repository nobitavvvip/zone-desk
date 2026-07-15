package controller

import (
	"zonedesk/internal/config"
	"zonedesk/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, shortcutSvc *service.ShortcutService, settingsSvc *service.SettingsService) *gin.Engine {
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

	containers := r.Group("/api/containers")
	{
		containers.GET("/list", ContainerListHandler())
		containers.POST("/start", ContainerStartHandler())
		containers.POST("/stop", ContainerStopHandler())
		containers.POST("/restart", ContainerRestartHandler())
		containers.POST("/remove", ContainerRemoveHandler())
		containers.GET("/logs", ContainerLogsHandler())
		containers.GET("/inspect", ContainerInspectHandler())
		containers.GET("/stats", ContainerStatsHandler())
		containers.POST("/prune", ContainerPruneHandler())
		containers.GET("/images", ContainerImagesHandler())
		containers.POST("/pull", ContainerPullImageHandler())
		containers.POST("/rmi", ContainerRemoveImageHandler())
		containers.POST("/create", ContainerCreateHandler())
		containers.POST("/exec", ContainerExecHandler())
	}

	compose := r.Group("/api/compose")
	{
		compose.GET("/list", ComposeListHandler())
		compose.POST("/up", ComposeUpHandler())
		compose.POST("/down", ComposeDownHandler())
		compose.POST("/start", ComposeStartHandler())
		compose.POST("/stop", ComposeStopHandler())
		compose.POST("/restart", ComposeRestartHandler())
		compose.POST("/create", ComposeCreateHandler())
		compose.GET("/read", ComposeReadHandler())
		compose.PUT("/update", ComposeUpdateHandler())
		compose.POST("/delete", ComposeDeleteHandler())
		compose.GET("/logs", ComposeLogsHandler())
		compose.GET("/ps", ComposePsHandler())
	}

	shortcuts := r.Group("/api/shortcuts")
	{
		shortcuts.GET("", ShortcutListHandler(shortcutSvc))
		shortcuts.POST("", ShortcutCreateHandler(shortcutSvc))
		shortcuts.DELETE("/:id", ShortcutDeleteHandler(shortcutSvc))
	}

	settings := r.Group("/api/settings")
	{
		settings.GET("", SettingsGetHandler(settingsSvc))
		settings.PUT("", SettingsSaveHandler(settingsSvc))
	}

	r.Static("/assets", cfg.Storage.WebDir+"/assets")
	r.NoRoute(func(c *gin.Context) {
		c.File(cfg.Storage.WebDir + "/index.html")
	})

	return r
}
