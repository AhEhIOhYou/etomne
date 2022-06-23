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

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	modelsList, err := modelModels.GetAllModels()
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "models.html", gin.H{
		"title": "Models",
		"data":  modelsList,
	})
}

func Model(c *gin.Context) {
	var req getModelRequest
	if err := c.ShouldBindUri(&req); err != nil {
		log.Fatal(err)
		return
	}

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	m, err := modelModels.GetModelById(req.ID)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, m)
}

func Upload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "Upload!",
	})
}

func UploadModel(c *gin.Context) {

	model := entities.Model3d{
		Name:        c.PostForm("name"),
		CreateDate:  time.Now().Format(time.RFC3339),
		Description: c.PostForm("descr"),
	}

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

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	model.FileId, err = models.CreateFile(&entities.File{Path: path}, server.Connect())
	if err != nil {
		return
	}

	modelModels.CreateModel(model)

	c.Redirect(http.StatusFound, "/models")
}

func Delete(c *gin.Context) {
	idString := c.Query("id")
	id, _ := strconv.Atoi(idString)

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	modelModels.DeleteModel(id)

	c.Redirect(http.StatusFound, "/models")
}
