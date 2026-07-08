package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	App struct {
		Root string `yaml:"root"`
	} `yaml:"app"`

	FileManager struct {
		DefaultPath    string `yaml:"defaultPath"`
		AllowRoot      bool   `yaml:"allowRoot"`
		MaxPreviewSize int64  `yaml:"maxPreviewSize"`
		MaxUploadSize  int64  `yaml:"maxUploadSize"`
	} `yaml:"fileManager"`

	Storage struct {
		DataDir       string `yaml:"dataDir"`
		LogDir        string `yaml:"logDir"`
		CacheDir      string `yaml:"cacheDir"`
		WebDir        string `yaml:"webDir"`
		ShortcutsFile string `yaml:"shortcutsFile"`
		SettingsFile  string `yaml:"settingsFile"`
	} `yaml:"storage"`

	Security struct {
		EnableLogin              bool `yaml:"enableLogin"`
		AllowDangerousOperations bool `yaml:"allowDangerousOperations"`
	} `yaml:"security"`

	UI struct {
		Theme       string `yaml:"theme"`
		DesktopMode bool   `yaml:"desktopMode"`
	} `yaml:"ui"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
