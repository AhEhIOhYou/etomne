package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"etomne/app/entities"
	"etomne/app/models"
	"etomne/app/server"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type getModelRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func Models(c *gin.Context) {
	modelsList, err := models.GetAllModels(server.Connect())
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, modelsList)
}

func Model(c *gin.Context) {
	var req getModelRequest
	if err := c.ShouldBindUri(&req); err != nil {
		log.Fatal(err)
		return
	}
	m, err := models.GetModelById(req.ID, server.Connect())
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, m)
}

func Up(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "Upload!",
	})
}

func Upload(c *gin.Context) {
	var err error

	c.Request.ParseMultipartForm(32 << 20)
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	h := md5.New()
	tmpFileName :=
		strconv.FormatInt(time.Now().Unix(), 10) +
			file.Filename +
			strconv.FormatInt(time.Now().UnixNano()%0x100000, 10)
	h.Write([]byte(tmpFileName))
	newFileName := hex.EncodeToString(h.Sum(nil))
	fileExtension := filepath.Ext(file.Filename)

	path := "upload/" + newFileName + fileExtension

	if err := c.SaveUploadedFile(file, path); err != nil {
		log.Fatal(err)
	}

	a, err := models.CreateFile(&entities.File{Path: path}, server.Connect())
	if err != nil {
		return
	}

	log.Println(a)

	c.JSON(http.StatusOK, gin.H{
		"Chisa": "nice body",
	})
}
