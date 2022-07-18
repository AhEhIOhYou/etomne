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
	fm       filemanager.ManagerFileInterface
	rd       auth.AuthInterface
	tk       auth.TokenInterface
}

func NewModel(mApp application.ModelAppInterface, uApp application.UserAppInterface, fApp application.FileAppInterface, fm filemanager.ManagerFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Model {
	return &Model{
		modelApp: mApp,
		userApp:  uApp,
		fileApp:  fApp,
		fm:       fm,
		rd:       rd,
		tk:       tk,
	}
}

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

	//var saveModelErr = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "invalid_json")
		return
	}

	emptyModel := entities.Model{
		Title:       title,
		Description: description,
	}

	saveModelErr := emptyModel.Validate("")
	if len(saveModelErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveModelErr)
		return
	}

	var Model = entities.Model{
		UserID:      userId,
		Title:       title,
		Description: description,
	}

	Model.Prepare()

	saveModel, saveErr := m.modelApp.SaveModel(&Model)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}

	files := c.Request.MultipartForm.File["attachments"]

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

		//Добавление файла в бд files
		File := entities.File{
			Title: file.Filename,
			Url:   url,
		}

		saveFile, saveModelErr := m.fileApp.SaveFile(&File)
		if len(saveModelErr) > 0 {
			c.JSON(http.StatusUnprocessableEntity, saveModelErr)
			return
		}

		//Добавление файла и модельки в бд model_files
		ModelFile := entities.ModelFile{
			ModelId: saveModel.ID,
			FileId:  saveFile.ID,
		}

		_, saveModelErr = m.modelApp.AddModelFile(&ModelFile)
		if len(saveModelErr) > 0 {
			c.JSON(http.StatusUnprocessableEntity, saveModelErr)
			return
		}
	}

	c.JSON(http.StatusCreated, saveModel)
}

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

	emptyModel := entities.Model{
		Title:       title,
		Description: description,
	}

	updateModelErr = emptyModel.Validate("update")
	if len(updateModelErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateModelErr)
		return
	}

	model.Title = title
	model.Description = description
	model.UpdatedAt = time.Now()

	model.BeforeSave()

	updatedModel, dbUpdateErr := m.modelApp.UpdateModel(model)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedModel)
}

func (m *Model) GetAllModel(c *gin.Context) {
	allModels, err := m.modelApp.GetAllModel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allModels)
}

// GetModel godoc
// @Summary      Get model
// @Description  Get model by ID
// @Tags         Models
// @Param        id   path      string  true  "Model ID"
// @Success      200  {object}  entities.Model
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
