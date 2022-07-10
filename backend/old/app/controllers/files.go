package controllers

import (
	"etomne/backend/old/app/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	MaxUploadSize = 200 >> 32
)

var AcceptableFormats = map[string]interface{}{
	"glb":  nil,
	"jpeg": nil,
	"jpg":  nil,
	"png":  nil,
}

func Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	file, fileHandler, err := c.Request.FormFile("file")
	if err != nil {
		server.NewError(c, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	buffer := make([]byte, fileHandler.Size)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)

	if _, ex := AcceptableFormats[fileType]; !ex {
		server.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

}
