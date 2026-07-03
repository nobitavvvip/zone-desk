package service

import (
	"zonedesk/internal/model"
	"zonedesk/internal/repository"
)

type DesktopService struct {
	store    *repository.JSONStore
	filePath string
}

func NewDesktopService(store *repository.JSONStore, filePath string) *DesktopService {
	return &DesktopService{
		store:    store,
		filePath: filePath,
	}
}

func (s *DesktopService) GetApps() ([]model.DesktopApp, error) {
	var apps []model.DesktopApp
	if err := s.store.Read(s.filePath, &apps); err != nil {
		return []model.DesktopApp{}, nil
	}
	return apps, nil
}

func (s *DesktopService) SaveApps(apps []model.DesktopApp) error {
	return s.store.Write(s.filePath, apps)
}
