package fileupload

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type fileUpload struct{}

func (fu *fileUpload) UploadFile(file *multipart.FileHeader) (string, error) {

	h := md5.New()
	tmpFileName :=
		strconv.FormatInt(time.Now().Unix(), 10) +
			file.Filename +
			strconv.FormatInt(time.Now().UnixNano()%0x100000, 10)
	h.Write([]byte(tmpFileName))

	newFileName := hex.EncodeToString(h.Sum(nil))
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

func NewFileUpload() *fileUpload {
	return &fileUpload{}
}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

var _ UploadFileInterface = &fileUpload{}
