package app

import (
	"zonedesk/internal/config"
	"zonedesk/internal/controller"
	"zonedesk/internal/repository"
	"zonedesk/internal/service"

	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config) *gin.Engine {
	store := repository.NewJSONStore()
	shortcutSvc := service.NewShortcutService(store, cfg.Storage.ShortcutsFile)
	settingsSvc := service.NewSettingsService(store, cfg.Storage.SettingsFile)

	r := controller.NewRouter(cfg, shortcutSvc, settingsSvc)
	return r
}
