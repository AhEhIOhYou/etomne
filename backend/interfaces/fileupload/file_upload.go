package fileupload

import (
	"errors"
	"fmt"
	"mime/multipart"
)

type fileUpload struct{}

func (fu *fileUpload) UploadFile(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer f.Close()

	size := file.Size

	fmt.Println("the size: ", size)

	// начать делать
	buffer := make([]byte, size)
	f.Read(buffer)
	//fileType := http.DetectContentType(buffer)
	filePath := FormatFile(file.Filename)

	return filePath, nil
}

func NewFileUpload() *fileUpload {
	return &fileUpload{}
}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

var _ UploadFileInterface = &fileUpload{}
