package controllers

import (
	"etomne/app/models"
	"etomne/app/server"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	//TODO сделать запись в бд + сгенерить имя файла

	if err := c.SaveUploadedFile(file, "upload/"+file.Filename); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"Chisa": "nice body",
	})
}
