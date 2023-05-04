package filemanager

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/security"
)

type fileManager struct{}

var _ ManagerFileInterface = &fileManager{}

type ManagerFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
	DeleteFile(path string) error
}

func NewFileUpload() *fileManager {
	return &fileManager{}
}

func (fu *fileManager) UploadFile(file *multipart.FileHeader) (string, error) {
	newFileName := security.CreateName(file.Filename)
	fileExtension := filepath.Ext(file.Filename)
	uploadDir := os.Getenv("UPLOAD_DIR")
	path := uploadDir + "/" + newFileName + fileExtension

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	return path, nil
}

func (fu *fileManager) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
