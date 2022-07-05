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

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func GetModels(c *gin.Context) {

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	modelsList, err := modelModels.GetAllModels()
	if err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"response": modelsList,
	})
}

func GetModel(c *gin.Context) {
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

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"response": m,
	})
}

func CreateModel(c *gin.Context) {

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

	lastId, _ := modelModels.CreateModel(model)

	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			"success": lastId,
		},
	})
}

func DeleteModel(c *gin.Context) {

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
	models.DeleteFile(m.FileId, server.Connect())
	lastId, _ := modelModels.DeleteModel(req.ID)

	c.JSON(http.StatusOK, gin.H{
		"response": lastId,
	})
}

func EditModel(c *gin.Context) {

	var m entities.Model3d
	c.ShouldBindJSON(&m)
	var arr map[string]any
	arr = make(map[string]any)
	arr["id"] = m.Id
	arr["name"] = m.Name

	log.Println(arr)

	//model := entities.Model3d{
	//	Name:        c.PostForm("name"),
	//	CreateDate:  c.uri
	//	Description: c.PostForm("descr"),
	//}
	//
	//c.Request.ParseMultipartForm(32 << 20)
	//file, err := c.FormFile("file")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//h := md5.New()
	//tmpFileName :=
	//	strconv.FormatInt(time.Now().Unix(), 10) +
	//		file.Filename +
	//		strconv.FormatInt(time.Now().UnixNano()%0x100000, 10)
	//h.Write([]byte(tmpFileName))
	//newFileName := hex.EncodeToString(h.Sum(nil))
	//fileExtension := filepath.Ext(file.Filename)
	//
	//path := "upload/" + newFileName + fileExtension
	//
	//if err := c.SaveUploadedFile(file, path); err != nil {
	//	log.Fatal(err)
	//}
	//
	//modelModels := models.Model3dModel{
	//	Db: server.Connect(),
	//}
	//
	//model.FileId, err = models.CreateFile(&entities.File{Path: path}, server.Connect())
	//if err != nil {
	//	return
	//}
	//
	//lastId, _ := modelModels.CreateModel(model)
	//
	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			"success": arr,
		},
	})
}
