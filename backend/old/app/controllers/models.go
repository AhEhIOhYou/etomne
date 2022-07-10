package controllers

import (
	models2 "etomne/backend/old/app/models"
	"etomne/backend/old/app/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getModelRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func GetModels(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	modelModels := models2.Model3dModel{
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
		server.NewResponse(c, http.StatusOK, gin.H{
			"count": len(modelsList),
			"items": modelsList,
		})
	}
}

func GetModel(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var req getModelRequest

	if err := c.ShouldBindUri(&req); err != nil {
		server.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	modelModels := models2.Model3dModel{
		Db: server.Connect(),
	}

	m, err := modelModels.GetModelById(req.ID)
	if err != nil {
		server.NewError(c, http.StatusBadRequest, err.Error())
		return
	}

	if m == "" {
		server.NewError(c, http.StatusNotFound, "Empty")
	} else {
		server.NewResponse(c, http.StatusOK, gin.H{
			"model": m,
		})
	}
}

func CreateModel(c *gin.Context) {

}

func DeleteModel(c *gin.Context) {

}

func EditModel(c *gin.Context) {

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
}
