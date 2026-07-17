package controller

import (
	"os"
	"path/filepath"
	"strings"

	"zonedesk/internal/config"
	"zonedesk/internal/service"
	"zonedesk/pkg/util"

	"github.com/gin-gonic/gin"
)

const defaultNotesDir = "./data/notes"

func getNotesRoot(cfg *config.Config) (string, error) {
	notesDir := cfg.Storage.NotesDir
	if notesDir == "" {
		notesDir = defaultNotesDir
	}
	abs, err := filepath.Abs(notesDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return "", err
	}
	return abs, nil
}

func NotesRootHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"rootDir": root})
	}
}

func NotesListHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		path := c.DefaultQuery("path", root)
		abs, err := util.CleanAbsPath(path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(abs, root) {
			util.Fail(c, 403, "access denied")
			return
		}

		result, err := service.ListFiles(abs, "name", "asc")
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, result)
	}
}

func NotesReadHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		path := c.Query("path")
		abs, err := util.CleanAbsPath(path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(abs, root) {
			util.Fail(c, 403, "access denied")
			return
		}

		result, err := service.ReadFile(abs, 10*1024*1024)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}

		util.OK(c, result)
	}
}

func NotesWriteHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		abs, err := util.CleanAbsPath(req.Path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(abs, root) {
			util.Fail(c, 403, "access denied")
			return
		}

		if err := service.WriteFile(abs, req.Content); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		item, err := service.StatFile(abs)
		if err != nil {
			util.OK(c, gin.H{"path": abs})
			return
		}

		util.OK(c, item)
	}
}

func NotesCreateHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		var req struct {
			DirPath string `json:"dirPath"`
			Name    string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		dirAbs, err := util.CleanAbsPath(req.DirPath)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(dirAbs, root) {
			util.Fail(c, 403, "access denied")
			return
		}

		fileName := req.Name
		if !strings.HasSuffix(fileName, ".md") {
			fileName += ".md"
		}
		fullPath := filepath.Join(dirAbs, fileName)
		title := strings.TrimSuffix(req.Name, ".md")
		content := "# " + title + "\n\n"

		if err := service.WriteFile(fullPath, content); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{"path": fullPath, "name": fileName})
	}
}

func NotesMkdirHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		var req struct {
			DirPath string `json:"dirPath"`
			Name    string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		dirAbs, err := util.CleanAbsPath(req.DirPath)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(dirAbs, root) {
			util.Fail(c, 403, "access denied")
			return
		}

		fullPath := filepath.Join(dirAbs, req.Name)
		if err := service.Mkdir(fullPath); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{"path": fullPath})
	}
}

func NotesDeleteHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		var req struct {
			Path string `json:"path"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		abs, err := util.CleanAbsPath(req.Path)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(abs, root) {
			util.Fail(c, 403, "access denied")
			return
		}
		if abs == root {
			util.Fail(c, 403, "cannot delete root")
			return
		}

		if err := service.Delete(abs); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{"path": abs})
	}
}

func NotesRenameHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		var req struct {
			OldPath string `json:"oldPath"`
			NewPath string `json:"newPath"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		oldAbs, err := util.CleanAbsPath(req.OldPath)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		newAbs, err := util.CleanAbsPath(req.NewPath)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(oldAbs, root) || !strings.HasPrefix(newAbs, root) {
			util.Fail(c, 403, "access denied")
			return
		}

		if err := service.Rename(oldAbs, newAbs); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{"oldPath": oldAbs, "newPath": newAbs})
	}
}

func NotesMoveHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		root, err := getNotesRoot(cfg)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		var req struct {
			Source      string `json:"source"`
			Destination string `json:"destination"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}

		src, err := util.CleanAbsPath(req.Source)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		dst, err := util.CleanAbsPath(req.Destination)
		if err != nil {
			util.Fail(c, 400, err.Error())
			return
		}
		if !strings.HasPrefix(src, root) || !strings.HasPrefix(dst, root) {
			util.Fail(c, 403, "access denied")
			return
		}

		if err := service.Move(src, dst); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}

		util.OK(c, gin.H{"source": src, "destination": dst})
	}
}
