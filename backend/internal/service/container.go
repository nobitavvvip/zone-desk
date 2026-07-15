package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"zonedesk/internal/model"
)

func ContainerList(all bool) ([]model.Container, error) {
	args := []string{"ps", "--format", "{{json .}}"}
	if all {
		args = append(args, "-a")
	}
	out, err := runDocker(args...)
	if err != nil {
		return nil, err
	}
	return parseContainerList(out)
}

func ContainerStart(id string) error {
	_, err := runDocker("start", id)
	return err
}

func ContainerStop(id string) error {
	_, err := runDocker("stop", id)
	return err
}

func ContainerRestart(id string) error {
	_, err := runDocker("restart", id)
	return err
}

func ContainerRemove(id string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, id)
	_, err := runDocker(args...)
	return err
}

func ContainerLogs(id string, tail int) (string, error) {
	args := []string{"logs", "--tail", fmt.Sprintf("%d", tail), "--timestamps", id}
	return runDocker(args...)
}

func ContainerInspect(id string) (string, error) {
	return runDocker("inspect", id)
}

func ContainerStats(id string) (*model.ContainerStats, error) {
	out, err := runDocker("stats", "--no-stream", "--format", "{{json .}}", id)
	if err != nil {
		return nil, err
	}
	return parseContainerStats(out)
}

func ContainerPrune() error {
	_, err := runDocker("container", "prune", "-f")
	return err
}

func ComposeList() ([]model.ComposeProject, error) {
	out, err := runDocker("compose", "ls", "--format", "json")
	if err != nil {
		// fallback: try without --format json (older compose versions)
		out2, err2 := runDocker("compose", "ls")
		if err2 != nil {
			return nil, err
		}
		out = out2
	}
	return parseComposeList(out)
}

func ComposeUp(projectDir, file string, detached bool) error {
	args := composeArgs(projectDir, file, "up")
	if detached {
		args = append(args, "-d")
	}
	_, err := runDocker(args...)
	return err
}

func ComposeDown(projectDir, file string) error {
	args := composeArgs(projectDir, file, "down")
	_, err := runDocker(args...)
	return err
}

func ComposeStart(projectDir, file string) error {
	args := composeArgs(projectDir, file, "start")
	_, err := runDocker(args...)
	return err
}

func ComposeStop(projectDir, file string) error {
	args := composeArgs(projectDir, file, "stop")
	_, err := runDocker(args...)
	return err
}

func ComposeRestart(projectDir, file string) error {
	args := composeArgs(projectDir, file, "restart")
	_, err := runDocker(args...)
	return err
}

func ComposePs(projectDir, file string) (string, error) {
	args := composeArgs(projectDir, file, "ps")
	return runDocker(args...)
}

func ComposeLogs(projectDir, file string, tail int) (string, error) {
	args := composeArgs(projectDir, file, "logs")
	args = append(args, "--tail", fmt.Sprintf("%d", tail), "--timestamps")
	return runDocker(args...)
}

func ComposeCreate(projectDir, filename, content string, start bool) error {
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	filePath := projectDir + "/" + filename
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	if start {
		args := composeArgs(projectDir, filename, "up")
		args = append(args, "-d")
		_, err := runDocker(args...)
		return err
	}
	return nil
}

func ComposeReadFile(projectDir, file string) (string, error) {
	filePath := projectDir + "/" + file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	return string(data), nil
}

func ComposeUpdateFile(projectDir, file, content string) error {
	filePath := projectDir + "/" + file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func ComposeDeleteProject(projectDir, file string) error {
	filePath := projectDir + "/" + file
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove file: %w", err)
	}
	// also try to remove empty project dir
	os.Remove(projectDir)
	return nil
}

func composeArgs(projectDir, file, subcmd string) []string {
	args := []string{"compose"}
	if file != "" {
		args = append(args, "-f", file)
	}
	if projectDir != "" {
		args = append(args, "--project-directory", projectDir)
	}
	args = append(args, subcmd)
	return args
}

func runDocker(args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("docker %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func RunDocker(args ...string) error {
	_, err := runDocker(args...)
	return err
}

func RunDockerOutput(args ...string) (string, error) {
	return runDocker(args...)
}

func parseContainerList(out string) ([]model.Container, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	var containers []model.Container
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var raw struct {
			ID      string `json:"ID"`
			Names   string `json:"Names"`
			Image   string `json:"Image"`
			Status  string `json:"Status"`
			State   string `json:"State"`
			Ports   string `json:"Ports"`
			CreatedAt string `json:"CreatedAt"`
		}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}
		containers = append(containers, model.Container{
			ID:      raw.ID,
			Name:    strings.TrimPrefix(raw.Names, "/"),
			Image:   raw.Image,
			Status:  raw.Status,
			State:   raw.State,
			Ports:   raw.Ports,
			Created: raw.CreatedAt,
		})
	}
	if containers == nil {
		containers = []model.Container{}
	}
	return containers, nil
}

func parseContainerStats(out string) (*model.ContainerStats, error) {
	var raw struct {
		CPUPercent    string `json:"CPUPerc"`
		MemUsage      string `json:"MemUsage"`
		MemLimit      string `json:"MemLimit"`
		NetIO         string `json:"NetIO"`
		BlockIO       string `json:"BlockIO"`
		PIDs          string `json:"PIDs"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, fmt.Errorf("parse stats: %w", err)
	}
	stats := &model.ContainerStats{}
	parsePercent(raw.CPUPercent, &stats.CPUPercent)
	parseMemBytes(raw.MemUsage, &stats.MemoryUsageMB)
	parseMemBytes(raw.MemLimit, &stats.MemoryLimitMB)
	parseIOBytes(raw.NetIO, &stats.NetworkRxMB, &stats.NetworkTxMB)
	parseIOBytes(raw.BlockIO, &stats.BlockRxMB, &stats.BlockWxMB)
	fmt.Sscanf(raw.PIDs, "%d", &stats.PIDs)
	return stats, nil
}

func parsePercent(s string, v *float64) {
	s = strings.TrimSuffix(s, "%")
	fmt.Sscanf(s, "%f", v)
}

func parseMemBytes(s string, v *float64) {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "/"); idx >= 0 {
		s = strings.TrimSpace(s[:idx])
	}
	switch {
	case strings.HasSuffix(s, "GiB"):
		fmt.Sscanf(s, "%fGiB", v)
		*v *= 1024
	case strings.HasSuffix(s, "MiB"):
		fmt.Sscanf(s, "%fMiB", v)
	case strings.HasSuffix(s, "KiB"):
		fmt.Sscanf(s, "%fKiB", v)
		*v /= 1024
	default:
		fmt.Sscanf(s, "%f", v)
	}
}

func parseIOBytes(s string, rx, tx *float64) {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "/"); idx >= 0 {
		rxStr := strings.TrimSpace(s[:idx])
		txStr := strings.TrimSpace(s[idx+1:])
		parseMemBytes(rxStr, rx)
		parseMemBytes(txStr, tx)
	}
}

func parseComposeList(out string) ([]model.ComposeProject, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return []model.ComposeProject{}, nil
	}

	// Try JSON array first (docker compose ls --format json)
	var raw []struct {
		Name        string `json:"Name"`
		Status      string `json:"Status"`
		ConfigFiles string `json:"ConfigFiles"`
	}
	if json.Unmarshal([]byte(out), &raw) == nil {
		return mapComposeProjects(raw), nil
	}

	// Fallback: parse table format
	// NAME                STATUS              CONFIG FILES
	// myproject           running(2)          /path/to/compose.yaml
	lines := strings.Split(out, "\n")
	// skip header line
	startIdx := 0
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "NAME") {
			startIdx = i + 1
			break
		}
	}
	for _, line := range lines[startIdx:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// split on 3+ spaces to handle columns with spaces in values
		parts := splitOnSpaces(line, 3)
		if len(parts) >= 2 {
			p := model.ComposeProject{
				Name:   parts[0],
				Status: parts[1],
			}
			if len(parts) >= 3 {
				p.ConfigFiles = parts[2]
			}
			raw = append(raw, struct {
				Name        string `json:"Name"`
				Status      string `json:"Status"`
				ConfigFiles string `json:"ConfigFiles"`
			}{Name: p.Name, Status: p.Status, ConfigFiles: p.ConfigFiles})
		}
	}
	return mapComposeProjects(raw), nil
}

func splitOnSpaces(s string, n int) []string {
	var result []string
	start := 0
	for i := 0; i < len(s) && len(result) < n-1; i++ {
		if s[i] == ' ' && i+1 < len(s) && s[i+1] == ' ' && s[i+2] == ' ' {
			result = append(result, strings.TrimSpace(s[start:i]))
			for i < len(s) && s[i] == ' ' {
				i++
			}
			start = i
		}
	}
	result = append(result, strings.TrimSpace(s[start:]))
	return result
}

func mapComposeProjects(raw []struct {
	Name        string `json:"Name"`
	Status      string `json:"Status"`
	ConfigFiles string `json:"ConfigFiles"`
}) []model.ComposeProject {
	var projects []model.ComposeProject
	for _, p := range raw {
		// extract directory from the first config file path
		wd := ""
		for _, cf := range splitComposeFiles(p.ConfigFiles) {
			cf = strings.TrimSpace(cf)
			if cf != "" {
				wd = cf
				break
			}
		}
		// workingDir should be the directory containing the compose file
		if idx := strings.LastIndex(wd, "/"); idx >= 0 {
			wd = wd[:idx]
		} else if idx := strings.LastIndex(wd, "\\"); idx >= 0 {
			wd = wd[:idx]
		}
		projects = append(projects, model.ComposeProject{
			Name:        p.Name,
			Status:      p.Status,
			ConfigFiles: p.ConfigFiles,
			WorkingDir:  wd,
		})
	}
	if projects == nil {
		projects = []model.ComposeProject{}
	}
	return projects
}

func splitComposeFiles(s string) []string {
	// ConfigFiles can be comma-separated for multi-file projects
	return strings.Split(s, ",")
}
