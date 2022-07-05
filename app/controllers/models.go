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

// GetModels godoc
// @Summary      Get models
// @Description  Get list models
// @Tags         Models
// @Success      200  {object} 	entities.Model3d
// @Success      204  {string}  string    "Empty"
// @Failure      500  {object}  server.HTTPError
// @Router       /models [get]
func GetModels(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	modelsList, err := modelModels.GetAllModels()
	if err != nil {
		server.WriteLog(server.Error, err.Error())
		server.NewError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(modelsList) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"response": modelsList,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"response": "empty",
		})
	}
}

// GetModel godoc
// @Summary      Get model
// @Description  Get model by ID
// @Tags         Models
// @Param        id   path      string  true  "Model ID"
// @Success      200  {object}  entities.Model3d
// @Failure      400  {object}  server.HTTPError
// @Failure      400  {object}  server.HTTPError
// @Failure      404  {object}  server.HTTPError
// @Router       /models/{id} [get]
func GetModel(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var req getModelRequest

	if err := c.ShouldBindUri(&req); err != nil {
		server.WriteLog(server.Error, err.Error())
		server.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	m, err := modelModels.GetModelById(req.ID)
	if err != nil {
		server.WriteLog(server.Error, err.Error())
		server.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	if m == (entities.Model3d{}) {
		server.NewError(c, http.StatusNotFound, "Empty")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"response": m,
		})
	}
}

func CreateModel(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	model := entities.Model3d{
		Name:        c.PostForm("name"),
		CreateDate:  time.Now().Format(time.RFC3339),
		Description: c.PostForm("descr"),
	}

	c.Request.ParseMultipartForm(32 << 20)
	file, err := c.FormFile("file")
	if err != nil {
		server.WriteLog(server.Error, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"response": gin.H{
				"error": err.Error(),
			},
		})
		return
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
		server.WriteLog(server.Error, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"response": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	model.FileId, err = models.CreateFile(&entities.File{Path: path}, server.Connect())
	if err != nil {
		server.WriteLog(server.Error, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"response": gin.H{
				"error": err.Error(),
			},
		})
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
	c.Header("Content-Type", "application/json")

	var req getModelRequest
	if err := c.ShouldBindUri(&req); err != nil {
		server.WriteLog(server.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"response": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	modelModels := models.Model3dModel{
		Db: server.Connect(),
	}

	m, err := modelModels.GetModelById(req.ID)
	if err != nil {
		server.WriteLog(server.Error, err.Error())
		c.JSON(http.StatusOK, gin.H{
			"response": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	if m == (entities.Model3d{}) {
		errMsg := "model not found"
		server.WriteLog(server.Error, errMsg)
		c.JSON(http.StatusOK, gin.H{
			"response": gin.H{
				"error": errMsg,
			},
		})
		return
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
