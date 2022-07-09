package interfaces

import (
	"etomne/application"
	"etomne/domain/entities"
	"etomne/infrastructure/auth"
	"etomne/interfaces/fileupload"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

	var Model = entities.Model{
		UserID:      userId,
		Title:       title,
		Description: description,
		ModelFile:   uploadedFile,
	}

	saveModel, saveErr := m.modelApp.SaveModel(&Model)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, saveModel)
}
