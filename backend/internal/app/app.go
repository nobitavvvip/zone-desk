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

	r := controller.NewRouter(cfg, shortcutSvc)
	return r
}
