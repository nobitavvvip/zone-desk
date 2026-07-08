package controller

import (
	"zonedesk/internal/model"
	"zonedesk/internal/service"
	"zonedesk/pkg/util"

	"github.com/gin-gonic/gin"
)

func SettingsGetHandler(svc *service.SettingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		settings, err := svc.Get()
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, settings)
	}
}

func SettingsSaveHandler(svc *service.SettingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var settings model.DesktopSettings
		if err := c.ShouldBindJSON(&settings); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		if settings.Blur < 0 || settings.Blur > 20 {
			util.Fail(c, 400, "blur must be 0-20")
			return
		}
		if settings.Mask < 0 || settings.Mask > 1 {
			util.Fail(c, 400, "mask must be 0-1")
			return
		}

		if err := svc.Save(&settings); err != nil {
			util.Fail(c, 500, "save settings failed")
			return
		}

		util.OK(c, settings)
	}
}
