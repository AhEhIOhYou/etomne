package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"etomne/app/entities"
	models2 "etomne/old/app/models"
	server2 "etomne/old/app/server"
	"github.com/gin-gonic/gin"
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
// @Success      200  {object} 	entities.Model
// @Success      204  {string}  string    "Empty"
// @Failure      500  {object}  server.HTTPError
// @Router       /models [get]
func GetModels(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	modelModels := models2.Model3dModel{
		Db: server2.Connect(),
	}

	modelsList, err := modelModels.GetAllModels()
	if err != nil {
		server2.WriteLog(server2.Error, err.Error())
		server2.NewError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(modelsList) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"response": modelsList,
		})
	} else {
		server2.NewResponse(c, http.StatusOK, gin.H{
			"count": len(modelsList),
			"items": modelsList,
		})
	}
}

// GetModel godoc
// @Summary      Get model
// @Description  Get model by ID
// @Tags         Models
// @Param        id   path      string  true  "Model ID"
// @Success      200  {object}  entities.Model
// @Failure      400  {object}  server.HTTPError
// @Failure      400  {object}  server.HTTPError
// @Failure      404  {object}  server.HTTPError
// @Router       /models/{id} [get]
func GetModel(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var req getModelRequest

	if err := c.ShouldBindUri(&req); err != nil {
		server2.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	modelModels := models2.Model3dModel{
		Db: server2.Connect(),
	}

	m, err := modelModels.GetModelById(req.ID)
	if err != nil {
		server2.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	if m == (entities.Model3d{}) {
		server2.NewError(c, http.StatusNotFound, "Empty")
	} else {
		server2.NewResponse(c, http.StatusOK, gin.H{
			"model": m,
		})
	}
}

// CreateModel godoc
// @Summary      Create model
// @Description  Create model
// @Tags         Models
// @Param        name	formData      string	true	"Model Title"
// @Param        descr	formData      string	false	"Model Description"
// @Param        file   formData      file		true	"Model File"
// @Success      200  {object}  server.HTTPResponse
// @Failure      400  {object}  server.HTTPError
// @Failure      404  {object}  server.HTTPError
// @Router       /models/create [post]
func CreateModel(c *gin.Context) {

	var m entities.NewModel3d
	err := c.ShouldBindJSON(&m)
	if err != nil {
		server2.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	server2.NewResponse(c, http.StatusOK, gin.H{
		"data": m,
	})
	return

	model := entities.Model3d{
		Name:        c.PostForm("name"),
		CreateDate:  time.Now().Format(time.RFC3339),
		Description: c.PostForm("descr"),
	}

	c.Request.ParseMultipartForm(32 << 20)
	file, err := c.FormFile("file")
	if err != nil {
		server2.NewError(c, http.StatusInternalServerError, err.Error())
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
		server2.NewError(c, http.StatusInternalServerError, err.Error())
		return
	}

	modelModels := models2.Model3dModel{
		Db: server2.Connect(),
	}

	model.FileId, err = models2.CreateFile(&entities.File{Path: path}, server2.Connect())
	if err != nil {
		server2.NewError(c, http.StatusInternalServerError, err.Error())
		return
	}

	lastId, _ := modelModels.CreateModel(model)

	server2.NewResponse(c, http.StatusOK, gin.H{
		"create": lastId,
	})
}

// DeleteModel godoc
// @Summary      Delete model
// @Description  Delete model
// @Tags         Models
// @Param        id   path      string  true  "Model ID"
// @Success      200  {object}  server.HTTPResponse
// @Failure      400  {object}  server.HTTPError
// @Failure      404  {object}  server.HTTPError
// @Failure      500  {object}  server.HTTPError
// @Router       /models/{id} [delete]
func DeleteModel(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var req getModelRequest
	if err := c.ShouldBindUri(&req); err != nil {
		server2.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	modelModels := models2.Model3dModel{
		Db: server2.Connect(),
	}

	m, err := modelModels.GetModelById(req.ID)
	if err != nil {
		server2.NewError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if m == (entities.Model3d{}) {
		server2.NewError(c, http.StatusNotFound, "model not found")
		return
	}
	models2.DeleteFile(m.FileId, server2.Connect())
	rowsAffected, _ := modelModels.DeleteModel(req.ID)

	server2.NewResponse(c, http.StatusOK, gin.H{
		"rows affected": rowsAffected,
	})
}

// EditModel godoc
// @Summary      Edit model
// @Description  Edit model
// @Tags         Models
// @Param        name	formData      string	true	"Model Title"
// @Param        descr	formData      string	false	"Model Description"
// @Param        file   formData      file		true	"Model File"
// @Success      200  {object}  server.HTTPResponse
// @Failure      400  {object}  server.HTTPError
// @Failure      404  {object}  server.HTTPError
// @Router       /models/{id} [put]
func EditModel(c *gin.Context) {

	var m entities.Model3d
	err := c.ShouldBindJSON(&m)
	if err != nil {
		server2.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	//model := entities.Model{
	//	Title:        c.PostForm("name"),
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
			"success": m,
		},
	})
}
