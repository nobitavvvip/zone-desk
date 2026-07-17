package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"zonedesk/internal/model"
)

func ListFiles(path string, sortBy string, sortOrder string) (*model.FileListResult, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	items := make([]model.FileItem, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileType := "file"
		if info.IsDir() {
			fileType = "dir"
		} else if info.Mode()&os.ModeSymlink != 0 {
			fileType = "symlink"
		}

		items = append(items, model.FileItem{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			Type:    fileType,
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().UTC().Format(time.RFC3339),
			Hidden:  strings.HasPrefix(entry.Name(), "."),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			if items[i].Type == "dir" {
				return true
			}
			if items[j].Type == "dir" {
				return false
			}
		}

		var result bool
		switch sortBy {
		case "size":
			result = items[i].Size < items[j].Size
		case "type":
			result = items[i].Type < items[j].Type
		case "modified":
			result = items[i].ModTime < items[j].ModTime
		default:
			result = strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		}

		if sortOrder == "desc" {
			return !result
		}
		return result
	})

	var parent *string
	if path != "/" {
		p := filepath.Dir(path)
		parent = &p
	}

	return &model.FileListResult{
		Path:   path,
		Parent: parent,
		Items:  items,
	}, nil
}

func StatFile(path string) (*model.FileItem, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	fileType := "file"
	if info.IsDir() {
		fileType = "dir"
	} else if info.Mode()&os.ModeSymlink != 0 {
		fileType = "symlink"
	}

	return &model.FileItem{
		Name:    filepath.Base(path),
		Path:    path,
		Type:    fileType,
		Size:    info.Size(),
		Mode:    info.Mode().String(),
		ModTime: info.ModTime().UTC().Format(time.RFC3339),
		Hidden:  strings.HasPrefix(filepath.Base(path), "."),
	}, nil
}

func ReadFile(path string, maxSize int64) (*model.FileReadResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.Mode().IsRegular() {
		return nil, errors.New("not a regular file")
	}

	if info.Size() > maxSize {
		return nil, errors.New("file too large for preview")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	encoding := "utf-8"

	return &model.FileReadResult{
		Path:     path,
		Content:  content,
		Encoding: encoding,
	}, nil
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func Delete(path string) error {
	if path == "/" {
		return errors.New("cannot delete root path")
	}
	return os.RemoveAll(path)
}

func Upload(targetDir, filename string) (string, error) {
	targetPath := filepath.Join(targetDir, filename)

	if _, err := os.Stat(targetPath); err == nil {
		return "", errors.New("file already exists")
	}

	return targetPath, nil
}

var TextExtensions = map[string]bool{
	".txt": true, ".log": true, ".md": true, ".json": true,
	".yaml": true, ".yml": true, ".toml": true, ".ini": true,
	".conf": true, ".env": true, ".sh": true, ".service": true,
	".xml": true, ".csv": true,
}

func IsTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return TextExtensions[ext]
}

var ImageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".webp": true, ".bmp": true, ".svg": true,
}

func IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ImageExtensions[ext]
}

func Copy(source, destination string) error {
	srcInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDir(source, destination)
	}
	return copyFile(source, destination)
}

func copyFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyDir(source, destination string) error {
	srcInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(destination, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(source, entry.Name())
		dstPath := filepath.Join(destination, entry.Name())

		if err := Copy(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func Move(source, destination string) error {
	return os.Rename(source, destination)
}

func WriteFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}
