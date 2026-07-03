package controller

import (
	"zonedesk/internal/service"
	"zonedesk/pkg/util"

	"github.com/gin-gonic/gin"
)

func ShortcutListHandler(svc *service.ShortcutService) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := svc.List()
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, items)
	}
}

func ShortcutCreateHandler(svc *service.ShortcutService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
			Path string `json:"path"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		if req.Name == "" {
			util.Fail(c, 400, "name is required")
			return
		}

		cleanPath, err := util.CleanAbsPath(req.Path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		sc, err := svc.Add(req.Name, cleanPath)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, sc)
	}
}

func ShortcutDeleteHandler(svc *service.ShortcutService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			util.Fail(c, 400, "id is required")
			return
		}

		if err := svc.Delete(id); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, nil)
	}
}
