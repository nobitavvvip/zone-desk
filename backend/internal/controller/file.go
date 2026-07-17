package controller

import (
	"path/filepath"

	"zonedesk/internal/config"
	"zonedesk/internal/service"
	"zonedesk/pkg/util"

	"github.com/gin-gonic/gin"
)

func FileListHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path, err := util.CleanAbsPath(c.Query("path"))
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		if !cfg.FileManager.AllowRoot && util.IsRootPath(path) {
			util.Fail(c, 403, "root access is not allowed")
			return
		}

		sortBy := c.DefaultQuery("sortBy", "name")
		sortOrder := c.DefaultQuery("sortOrder", "asc")

		result, err := service.ListFiles(path, sortBy, sortOrder)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, result)
	}
}

func FileStatHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path, err := util.CleanAbsPath(c.Query("path"))
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		item, err := service.StatFile(path)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, item)
	}
}

func FileReadHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path, err := util.CleanAbsPath(c.Query("path"))
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		result, err := service.ReadFile(path, cfg.FileManager.MaxPreviewSize)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		util.OK(c, result)
	}
}

func FileDownloadHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path, err := util.CleanAbsPath(c.Query("path"))
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		c.FileAttachment(path, filepath.Base(path))
	}
}

func FileUploadHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetDir := c.PostForm("path")
		targetDir, err := util.CleanAbsPath(targetDir)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			util.Fail(c, 400, "file is required")
			return
		}

		savePath, err := service.Upload(targetDir, file.Filename)
		if err != nil {
			util.Fail(c, 409, err.Error())
			return
		}

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			util.Fail(c, 500, "failed to save file")
			return
		}

		util.OK(c, gin.H{
			"path": savePath,
		})
	}
}

func FileMkdirHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Path string `json:"path"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		path, err := util.CleanAbsPath(req.Path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		if err := service.Mkdir(path); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{"path": path})
	}
}

func FileRenameHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			OldPath string `json:"oldPath"`
			NewPath string `json:"newPath"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		oldPath, err := util.CleanAbsPath(req.OldPath)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		newPath, err := util.CleanAbsPath(req.NewPath)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		if err := service.Rename(oldPath, newPath); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{
			"oldPath": oldPath,
			"newPath": newPath,
		})
	}
}

func FileCopyHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Source      string `json:"source"`
			Destination string `json:"destination"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		source, err := util.CleanAbsPath(req.Source)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		destination, err := util.CleanAbsPath(req.Destination)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		if err := service.Copy(source, destination); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{
			"source":      source,
			"destination": destination,
		})
	}
}

func FileMoveHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Source      string `json:"source"`
			Destination string `json:"destination"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		source, err := util.CleanAbsPath(req.Source)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		destination, err := util.CleanAbsPath(req.Destination)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		if err := service.Move(source, destination); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{
			"source":      source,
			"destination": destination,
		})
	}
}

func FileDeleteHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Path string `json:"path"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		path, err := util.CleanAbsPath(req.Path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		if err := service.Delete(path); err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		util.OK(c, gin.H{"path": path})
	}
}

func FileWriteHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		path, err := util.CleanAbsPath(req.Path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		if err := service.WriteFile(path, req.Content); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		item, err := service.StatFile(path)
		if err != nil {
			util.OK(c, gin.H{"path": path})
			return
		}

		util.OK(c, item)
	}
}
