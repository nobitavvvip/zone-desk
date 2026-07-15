package controller

import (
	"strconv"

	"zonedesk/internal/service"
	"zonedesk/pkg/util"

	"github.com/gin-gonic/gin"
)

func ContainerListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		all := c.DefaultQuery("all", "false") == "true"
		containers, err := service.ContainerList(all)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, containers)
	}
}

func ContainerStartHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID string `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ContainerStart(req.ID); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"id": req.ID})
	}
}

func ContainerStopHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID string `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ContainerStop(req.ID); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"id": req.ID})
	}
}

func ContainerRestartHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID string `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ContainerRestart(req.ID); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"id": req.ID})
	}
}

func ContainerRemoveHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID    string `json:"id"`
			Force bool   `json:"force"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ContainerRemove(req.ID, req.Force); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"id": req.ID})
	}
}

func ContainerLogsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			util.Fail(c, 400, "id is required")
			return
		}
		tail := 100
		if t := c.Query("tail"); t != "" {
			if v, err := strconv.Atoi(t); err == nil && v > 0 {
				tail = v
			}
		}
		logs, err := service.ContainerLogs(id, tail)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"logs": logs})
	}
}

func ContainerInspectHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			util.Fail(c, 400, "id is required")
			return
		}
		info, err := service.ContainerInspect(id)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"inspect": info})
	}
}

func ContainerStatsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			util.Fail(c, 400, "id is required")
			return
		}
		stats, err := service.ContainerStats(id)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, stats)
	}
}

func ContainerPruneHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := service.ContainerPrune(); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"pruned": true})
	}
}

func ComposeListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		projects, err := service.ComposeList()
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, projects)
	}
}

func ComposeUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			File       string `json:"file"`
			Detached   bool   `json:"detached"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ComposeUp(req.ProjectDir, req.File, req.Detached); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"status": "up"})
	}
}

func ComposeDownHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			File       string `json:"file"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ComposeDown(req.ProjectDir, req.File); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"status": "down"})
	}
}

func ComposeStartHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			File       string `json:"file"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ComposeStart(req.ProjectDir, req.File); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"status": "started"})
	}
}

func ComposeStopHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			File       string `json:"file"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ComposeStop(req.ProjectDir, req.File); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"status": "stopped"})
	}
}

func ComposeRestartHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			File       string `json:"file"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.ComposeRestart(req.ProjectDir, req.File); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"status": "restarted"})
	}
}

func ComposeLogsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectDir := c.Query("projectDir")
		file := c.Query("file")
		tail := 100
		if t := c.Query("tail"); t != "" {
			if v, err := strconv.Atoi(t); err == nil && v > 0 {
				tail = v
			}
		}
		logs, err := service.ComposeLogs(projectDir, file, tail)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"logs": logs})
	}
}

func ComposeCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			Filename   string `json:"filename"`
			Content    string `json:"content"`
			Start      bool   `json:"start"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		if req.Filename == "" {
			req.Filename = "compose.yaml"
		}
		if req.ProjectDir == "" {
			req.ProjectDir = "."
		}
		if err := service.ComposeCreate(req.ProjectDir, req.Filename, req.Content, req.Start); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"projectDir": req.ProjectDir, "filename": req.Filename})
	}
}

func ComposeReadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectDir := c.Query("projectDir")
		file := c.Query("file")
		if file == "" {
			file = "compose.yaml"
		}
		content, err := service.ComposeReadFile(projectDir, file)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"content": content, "projectDir": projectDir, "filename": file})
	}
}

func ComposeUpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			Filename   string `json:"filename"`
			Content    string `json:"content"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		if req.Filename == "" {
			req.Filename = "compose.yaml"
		}
		if err := service.ComposeUpdateFile(req.ProjectDir, req.Filename, req.Content); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"projectDir": req.ProjectDir, "filename": req.Filename})
	}
}

func ComposeDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ProjectDir string `json:"projectDir"`
			Filename   string `json:"filename"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.Fail(c, 400, "invalid request")
			return
		}
		if req.Filename == "" {
			req.Filename = "compose.yaml"
		}
		if err := service.ComposeDeleteProject(req.ProjectDir, req.Filename); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"deleted": true})
	}
}

func ComposePsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectDir := c.Query("projectDir")
		file := c.Query("file")
		out, err := service.ComposePs(projectDir, file)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"output": out})
	}
}

func ContainerImagesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		out, err := service.RunDockerOutput("images", "--format", "{{json .}}")
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"output": out})
	}
}

func ContainerPullImageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Image string `json:"image"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Image == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		if err := service.RunDocker("pull", req.Image); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"image": req.Image})
	}
}

func ContainerRemoveImageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Image string `json:"image"`
			Force bool   `json:"force"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Image == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		args := []string{"rmi"}
		if req.Force {
			args = append(args, "-f")
		}
		args = append(args, req.Image)
		if err := service.RunDocker(args...); err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"image": req.Image})
	}
}

func ContainerCreateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name       string   `json:"name"`
			Image      string   `json:"image"`
			Ports      []string `json:"ports"`
			Env        []string `json:"env"`
			Volumes    []string `json:"volumes"`
			Network    string   `json:"network"`
			Entrypoint string   `json:"entrypoint"`
			Cmd        []string `json:"cmd"`
			Restart    string   `json:"restart"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Image == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		args := []string{"create", "--name", req.Name}
		for _, p := range req.Ports {
			if p != "" {
				args = append(args, "-p", p)
			}
		}
		for _, e := range req.Env {
			if e != "" {
				args = append(args, "-e", e)
			}
		}
		for _, v := range req.Volumes {
			if v != "" {
				args = append(args, "-v", v)
			}
		}
		if req.Network != "" {
			args = append(args, "--network", req.Network)
		}
		if req.Entrypoint != "" {
			args = append(args, "--entrypoint", req.Entrypoint)
		}
		if req.Restart != "" {
			args = append(args, "--restart", req.Restart)
		}
		args = append(args, req.Image)
		args = append(args, req.Cmd...)
		out, err := service.RunDockerOutput(args...)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"containerId": out})
	}
}

func ContainerExecHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID      string   `json:"id"`
			Cmd     []string `json:"cmd"`
			Detach  bool     `json:"detach"`
			Interactive bool `json:"interactive"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
			util.Fail(c, 400, "invalid request")
			return
		}
		args := []string{"exec"}
		if req.Detach {
			args = append(args, "-d")
		}
		if req.Interactive {
			args = append(args, "-it")
		}
		args = append(args, req.ID)
		args = append(args, req.Cmd...)
		out, err := service.RunDockerOutput(args...)
		if err != nil {
			util.Fail(c, 500, err.Error())
			return
		}
		util.OK(c, gin.H{"output": out})
	}
}
