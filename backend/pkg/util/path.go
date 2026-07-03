package util

import (
	"errors"
	"path/filepath"
	"strings"
)

func CleanAbsPath(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("path is empty")
	}

	clean := filepath.Clean(input)

	if !filepath.IsAbs(clean) {
		return "", errors.New("path must be absolute")
	}

	return clean, nil
}

func IsRootPath(path string) bool {
	return filepath.Clean(path) == "/"
}
