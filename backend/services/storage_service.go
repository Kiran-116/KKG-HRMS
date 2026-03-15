package services

import (
	"io"
	"os"
	"path/filepath"

	"hrms/config"
)

type StorageService interface {
	SaveFile(filename string, file io.Reader) (string, error)
	GetFile(filepath string) (io.ReadCloser, error)
	DeleteFile(filepath string) error
}

type LocalStorageService struct {
	basePath string
}

func NewLocalStorageService() StorageService {
	basePath := config.AppConfig.Storage.LocalPath
	if basePath == "" {
		basePath = "./uploads"
	}

	// Create directory if it doesn't exist
	os.MkdirAll(basePath, 0755)

	return &LocalStorageService{
		basePath: basePath,
	}
}

func (s *LocalStorageService) SaveFile(filename string, file io.Reader) (string, error) {
	filePath := filepath.Join(s.basePath, filename)
	
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filePath, nil
}

func (s *LocalStorageService) GetFile(filepath string) (io.ReadCloser, error) {
	return os.Open(filepath)
}

func (s *LocalStorageService) DeleteFile(filepath string) error {
	return os.Remove(filepath)
}
