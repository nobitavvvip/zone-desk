package model

type Container struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	State   string `json:"state"`
	Ports   string `json:"ports"`
	Created string `json:"created"`
}

type ContainerStats struct {
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryUsageMB float64 `json:"memoryUsageMB"`
	MemoryLimitMB float64 `json:"memoryLimitMB"`
	NetworkRxMB   float64 `json:"networkRxMB"`
	NetworkTxMB   float64 `json:"networkTxMB"`
	BlockRxMB     float64 `json:"blockRxMB"`
	BlockWxMB     float64 `json:"blockWxMB"`
	PIDs          int     `json:"pids"`
}

type ComposeProject struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	ConfigFiles string `json:"configFiles"`
	WorkingDir  string `json:"workingDir"`
}
