package service

import (
	"crypto/rand"
	"encoding/hex"

	"zonedesk/internal/model"
	"zonedesk/internal/repository"
)

type ShortcutService struct {
	store    *repository.JSONStore
	filePath string
}

func NewShortcutService(store *repository.JSONStore, filePath string) *ShortcutService {
	return &ShortcutService{
		store:    store,
		filePath: filePath,
	}
}

func (s *ShortcutService) List() ([]model.Shortcut, error) {
	var shortcuts []model.Shortcut
	if err := s.store.Read(s.filePath, &shortcuts); err != nil {
		return []model.Shortcut{}, nil
	}
	return shortcuts, nil
}

func (s *ShortcutService) Add(name, path string) (model.Shortcut, error) {
	shortcuts, err := s.List()
	if err != nil {
		return model.Shortcut{}, err
	}

	id, err := generateID()
	if err != nil {
		return model.Shortcut{}, err
	}

	sc := model.Shortcut{
		ID:   id,
		Name: name,
		Path: path,
	}

	shortcuts = append(shortcuts, sc)

	if err := s.store.Write(s.filePath, shortcuts); err != nil {
		return model.Shortcut{}, err
	}

	return sc, nil
}

func (s *ShortcutService) Delete(id string) error {
	shortcuts, err := s.List()
	if err != nil {
		return err
	}

	var updated []model.Shortcut
	for _, sc := range shortcuts {
		if sc.ID != id {
			updated = append(updated, sc)
		}
	}

	return s.store.Write(s.filePath, updated)
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
