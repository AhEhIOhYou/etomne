package interfaces

import (
	"fmt"
	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/filemanager"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Model struct {
	modelApp application.ModelAppInterface
	userApp  application.UserAppInterface
	fileApp  application.FileAppInterface
	comApp   application.CommentAppInterface
	fm       filemanager.ManagerFileInterface
	rd       auth.AuthInterface
	tk       auth.TokenInterface
}

func NewModel(mApp application.ModelAppInterface, uApp application.UserAppInterface, fApp application.FileAppInterface, comApp application.CommentAppInterface, fm filemanager.ManagerFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Model {
	return &Model{
		modelApp: mApp,
		userApp:  uApp,
		fileApp:  fApp,
		comApp:   comApp,
		fm:       fm,
		rd:       rd,
		tk:       tk,
	}
}

// SaveModel godoc
// @Summary      Save model
// @Tags         Model
// @Accept       mpfd
// @Produce      json
// @Param        title		   formData      string  true  "Model Title"
// @Param        description   formData      string  true  "Model Description"
// @Param        attachments   formData      file	 false "Model Files"		Format(binary)
// @Success      201  {object}  entities.Model
// @Failure      401  string  unauthorized
// @Failure      400  string  user not found, unauthorized
// @Failure      422  string  error
// @Failure      500  string  error
// @Router       /model [post]
// @Security	 bearerAuth
func (m *Model) SaveModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = m.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "invalid_json")
		return
	}

	var Model = entities.Model{
		UserID:      userId,
		Title:       title,
		Description: description,
	}

	Model.Prepare()

	saveModelErr := Model.Validate("")
	if len(saveModelErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveModelErr)
		return
	}

	saveModel, saveErr := m.modelApp.SaveModel(&Model)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}

	files := c.Request.MultipartForm.File["attachments"]
	if len(files) > 0 {
		for _, file := range files {

			//Загрузка файла на сервер
			url, err := m.fm.UploadFile(file)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, err.Error())
				return
			}
			if url == "" {
				c.JSON(http.StatusUnprocessableEntity, "something went wrong")
				return
			}

			File := entities.File{
				OwnerId: userId,
				Title:   file.Filename,
				Url:     url,
			}

			File.Prepare()
			saveFileErr := File.Validate("")
			if len(saveFileErr) > 0 {
				c.JSON(http.StatusUnprocessableEntity, saveFileErr)
				return
			}

			_, saveFileErr = m.modelApp.SaveModelFile(&File, Model.ID)
			if len(saveFileErr) > 0 {
				c.JSON(http.StatusUnprocessableEntity, saveFileErr)
				return
			}
		}
	}

	c.JSON(http.StatusCreated, saveModel)
}

// UpdateModel godoc
// @Summary      Update model
// @Tags         Model
// @Accept       mpfd
// @Produce      json
// @Param        id		   	   path      string  true  "Model ID"
// @Param        title		   formData  string  true  "Model Title"
// @Param        description   formData  string  true  "Model Description"
// @Success      201  {object}  entities.Model
// @Failure      401  string  unauthorized
// @Failure      400  string  invalid request
// @Failure      422  string  error
// @Failure      500  string  error
// @Router       /model/{id} [put]
// @Security	 bearerAuth
func (m *Model) UpdateModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := m.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	modelId, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	model, err := m.modelApp.GetModel(modelId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	if user.ID != model.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this model")
		return
	}

	var updateModelErr = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")

	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}

	model.Title = title
	model.Description = description
	model.UpdatedAt = time.Now()

	model.BeforeUpdate()

	updateModelErr = model.Validate("update")
	if len(updateModelErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateModelErr)
		return
	}

	updatedModel, dbUpdateErr := m.modelApp.UpdateModel(model)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}

	c.JSON(http.StatusOK, updatedModel)
}

// GetAllModel godoc
// @Summary      Get all models
// @Tags         Model
// @Produce      json
// @Success      200  {array}  entities.Model
// @Failure      401  string  unauthorized
// @Failure      400  string  invalid request
// @Failure      422  string  error
// @Failure      500  string  error
// @Router       /model [get]
func (m *Model) GetAllModel(c *gin.Context) {

	limit, err := strconv.Atoi(c.Query("_limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	page, err := strconv.Atoi(c.Query("_page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	allModels, err := m.modelApp.GetAllModel(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allModels)
}

// GetModel godoc
// @Summary      Get model
// @Tags         Model
// @Param        id   path      string  true  "Model ID"
// @Success      200  {object}  entities.Model
// @Failure      401  string  unauthorized
// @Failure      400  string  invalid request
// @Failure      422  string  error
// @Failure      500  string  error
// @Router       /model/{id} [get]
func (m *Model) GetModel(c *gin.Context) {
	modelId, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	model, err := m.modelApp.GetModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := m.userApp.GetUser(model.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	files, err := m.modelApp.GetFilesByModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	Model := map[string]interface{}{
		"model":  model,
		"author": user.PublicUser(),
		"files":  files,
	}
	c.JSON(http.StatusOK, Model)
}

// DeleteModel godoc
// @Summary      Delete model
// @Tags         Model
// @Param        id   path      string  true  "Model ID"
// @Success      200  {string} string  model deleted
// @Failure      401  string  unauthorized
// @Failure      400  string  invalid request
// @Failure      422  string  error
// @Failure      500  string  error
// @Router       /model/{id} [delete]
// @Security	 bearerAuth
func (m *Model) DeleteModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	modelId, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = m.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	isAvaliable, err := m.modelApp.CheckAvailabilityModel(modelId, metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !isAvaliable {
		c.JSON(http.StatusInternalServerError, "model is unavailable")
		return
	}

	//Очистка комментариев
	err = m.comApp.DeleteCommentsByModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//Получение всех файлов модельки
	files, err := m.modelApp.GetFilesByModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//Отвязывание файлов от модели
	err = m.modelApp.DeleteAllModelFiles(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//Удаление файлов с бд и хранилища
	for _, file := range files {
		var deleteFileErr error
		deleteFileErr = m.fileApp.DeleteFile(file.ID)
		if deleteFileErr != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		deleteFileErr = m.fm.DeleteFile(file.Url)
		if deleteFileErr != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	//Удаление модели
	err = m.modelApp.DeleteModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "model deleted")
}

// SaveModelFile doc
// @Summary		Save model file
// @Tags		File
// @Accept		mpfd
// @Produce		json
// @Param		model_id  formData  string  false  "Model ID"
// @Param		file      formData  file    true  "File"
// @Success		201  {object}  entities.File
// @Failure     401  string  unauthorized
// @Failure     400  string  error
// @Failure     422  string  error
// @Failure     500  string  error
// @Router		/model/addfile/ [post]
// @Security	bearerAuth
func (m *Model) SaveModelFile(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = m.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	modelId, err := strconv.ParseUint(c.PostForm("model_id"), 10, 64)
	if err != nil || modelId == 0 {
		c.JSON(http.StatusBadRequest, "invalid query")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	url, err := m.fm.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if url == "" {
		c.JSON(http.StatusUnprocessableEntity, "something went wrong")
		return
	}

	File := entities.File{
		OwnerId: userId,
		Title:   file.Filename,
		Url:     url,
	}

	File.Prepare()
	saveFileErr := File.Validate("")
	if len(saveFileErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveFileErr)
		return
	}

	modelFile, saveFileErr := m.modelApp.SaveModelFile(&File, modelId)
	if len(saveFileErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveFileErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file":       File,
		"model_file": modelFile,
	})
}
