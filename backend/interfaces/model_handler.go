package interfaces

import (
	"etomne/backend/application"
	"etomne/backend/domain/entities"
	"etomne/backend/infrastructure/auth"
	"etomne/backend/interfaces/fileupload"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Model struct {
	modelApp   application.ModelAppInterface
	userApp    application.UserAppInterface
	fileUpload fileupload.UploadFileInterface
	rd         auth.AuthInterface
	tk         auth.TokenInterface
}

func NewModel(mApp application.ModelAppInterface, uApp application.UserAppInterface, fd fileupload.UploadFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Model {
	return &Model{
		modelApp:   mApp,
		userApp:    uApp,
		fileUpload: fd,
		rd:         rd,
		tk:         tk,
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

	var saveModelErr = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}

	emptyModel := entities.Model{
		Title:       title,
		Description: description,
	}

	saveModelErr = emptyModel.Validate("")
	if len(saveModelErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveModelErr)
		return
	}

	// В разработке
	/*
		file, err := c.FormFile("model_file")
		if err != nil {
			saveModelErr["invalid_file"] = "a valid file is required"
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}

		uploadedFile, err := m.fileUpload.UploadFile(file)
		if err != nil {
			saveModelErr["upload_err"] = err.Error()
			c.JSON(http.StatusUnprocessableEntity, saveModelErr)
			return
		}
	*/

	var Model = entities.Model{
		UserID:      userId,
		Title:       title,
		Description: description,
	}

	saveModel, saveErr := m.modelApp.SaveModel(&Model)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
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
		c.JSON(http.StatusUnauthorized, "unauthorized")
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

	// В разработке
	//file, _ := c.FormFile("model_file")

	/*
		if file != nil {
			model.ModelFile, err = m.fileUpload.UploadFile(file)

			model.ModelFile = os.Getenv("DO_SPACES_URL") + model.ModelFile
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"upload_err": err.Error(),
				})
				return
			}
		}
	*/

	model.Title = title
	model.Description = description
	model.UpdatedAt = time.Now()
	updatedFood, dbUpdateErr := m.modelApp.UpdateModel(model)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedFood)
}

func (m *Model) GetAllModel(c *gin.Context) {
	allModels, err := m.modelApp.GetAllModel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allModels)
}

// GetModelAndAuthor godoc
// @Summary      Get model and author
// @Description  Get model and author by ID model
// @Tags         Models
// @Param        id   path      string  true  "Model ID"
// @Success      200  {object}  entities.Model
// @Router       /model/{id} [get]
func (m *Model) GetModelAndAuthor(c *gin.Context) {
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
	foodAndUser := map[string]interface{}{
		"model":  model,
		"author": user.PublicUser(),
	}
	c.JSON(http.StatusOK, foodAndUser)
}

func (m *Model) DeleteModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
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
	err = m.modelApp.DeleteModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "model deleted")
}
