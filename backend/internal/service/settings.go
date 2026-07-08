package service

import (
	"zonedesk/internal/model"
	"zonedesk/internal/repository"
)

type SettingsService struct {
	store    *repository.JSONStore
	filePath string
}

func NewSettingsService(store *repository.JSONStore, filePath string) *SettingsService {
	return &SettingsService{
		store:    store,
		filePath: filePath,
	}
}

func (s *SettingsService) Get() (*model.DesktopSettings, error) {
	var settings model.DesktopSettings
	if err := s.store.Read(s.filePath, &settings); err != nil {
		return &model.DesktopSettings{
			Wallpaper: "none",
			Blur:      0,
			Mask:      0.35,
			Accent:    "#3b82f6",
		}, nil
	}
	return &settings, nil
}

func (s *SettingsService) Save(settings *model.DesktopSettings) error {
	return s.store.Write(s.filePath, settings)
}
