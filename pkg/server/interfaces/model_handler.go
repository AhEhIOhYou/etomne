package interfaces

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/AhEhIOhYou/etomne/pkg/server/application"
	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/utils"
	"github.com/AhEhIOhYou/etomne/pkg/server/interfaces/filemanager"
	"github.com/gin-gonic/gin"
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

//	@Summary	Save model
//	@Tags		model
//	@Produce	json
//	@Param		data	body		entities.ModelRequest	true	"Model data"
//	@Success	201		{object}	entities.Model
//	@Failure	400		string		string
//	@Failure	401		string		string
//	@Failure	500		string		string
//	@Router		/model [post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) SaveModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	_, err = m.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	var modelReq entities.ModelRequest

	if err := c.ShouldBindJSON(&modelReq); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	validateErr := modelReq.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	newModel := modelReq.NewModel()
	newModel.Prepare()

	savedModel, err := m.modelApp.SaveModel(newModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusCreated, savedModel)
}

//	@Summary	Update model
//	@Tags		model
//	@Produce	json
//	@Param		model_id	path		int						true	"Model ID"
//	@Param		data		body		entities.ModelRequest	true	"Model updated data"
//	@Success	200			{object}	entities.Model
//	@Failure	400			string		string
//	@Failure	401			string		string
//	@Failure	500			string		string
//	@Router		/model/update/{model_id} [post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) UpdateModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	user, err := m.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelID, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	updatableModel, err := m.modelApp.GetModel(modelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	if !utils.AccessVerification(updatableModel.UserID, user, false) {
		c.JSON(http.StatusUnauthorized, constants.NotEnoughRights)
		return
	}

	var modelReq entities.ModelRequest

	if err := c.ShouldBindJSON(&modelReq); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	if (modelReq == entities.ModelRequest{}) {
		c.JSON(http.StatusOK, updatableModel)
		return
	}

	if len(modelReq.Title) != 0 {
		updatableModel.Title = modelReq.Title
	}
	if len(modelReq.Description) != 0 {
		updatableModel.Description = modelReq.Description
	}

	updatableModel.BeforeUpdate()

	validateErr := updatableModel.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	updatedModel, err := m.modelApp.UpdateModel(updatableModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, updatedModel)
}

//	@Summary	Get a list of models with the specified quantity and position
//	@Tags		model
//	@Param		_page	query	int	false	"Query page param"
//	@Param		_limit	query	int	false	"Query limit param"
//	@Success	201		{Array}	[]entities.ModelData
//	@Failure	400		string	string
//	@Failure	401		string	string
//	@Failure	500		string	string
//	@Router		/model   [put]
func (m *Model) GetModelList(c *gin.Context) {

	limit, err := strconv.Atoi(c.Query("_limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	page, err := strconv.Atoi(c.Query("_page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	var readyModels []entities.ModelData

	rawModels, err := m.modelApp.GetAllModel(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	for _, model := range rawModels {

		user, err := m.userApp.GetUser(model.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}

		files, err := m.modelApp.GetFilesByModel(model.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}

		readyModels = append(readyModels, entities.ModelData{
			Model: model,
			User:  *user.PublicUser(),
			Files: utils.SortFiles(files),
		})
	}

	c.JSON(http.StatusOK, readyModels)
}

//	@Summary	Get model by ID
//	@Tags		model
//	@Param		model_id	path		int	true	"Model ID"
//	@Success	201			{object}	entities.ModelData
//	@Failure	400			string		string
//	@Failure	401			string		string
//	@Failure	500			string		string
//	@Router		/model/{model_id}  [get]
func (m *Model) GetModel(c *gin.Context) {
	modelId, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	model, err := m.modelApp.GetModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	user, err := m.userApp.GetUser(model.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	files, err := m.modelApp.GetFilesByModel(modelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, &entities.ModelData{
		Model: *model,
		User:  *user.PublicUser(),
		Files: utils.SortFiles(files),
	})
}

//	@Summary	Delete model by ID
//	@Tags		model
//	@Param		model_id	path	int	true	"Model ID"
//	@Success	200			string	string
//	@Failure	400			string	string
//	@Failure	401			string	string
//	@Failure	500			string	string
//	@Router		/model/{model_id} [delete]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) DeleteModel(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelID, err := strconv.ParseUint(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	user, err := m.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	model, err := m.modelApp.GetModel(modelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	if !utils.AccessVerification(model.UserID, user, false) {
		c.JSON(http.StatusUnauthorized, constants.NotEnoughRights)
		return
	}

	files, err := m.modelApp.GetFilesByModel(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	// Delete files from storage
	for _, file := range files {
		err = m.fm.DeleteFile(file.Url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
			return
		}
	}

	err = m.modelApp.DeleteModel(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, constants.DeletedSuccessful)
}

//	@Summary	Save model file
//	@Tags		file
//	@Produce	json
//	@Param		model_id	path		int		true	"Model ID"
//	@Param		file		formData	file	true	"Body with files"
//	@Success	201			{object}	entities.File
//	@Failure	400			string		string
//	@Failure	401			string		string
//	@Failure	500			string		string
//	@Router		/model/{model_id}/addfile     [post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (m *Model) SaveModelFile(c *gin.Context) {
	metadata, err := m.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := m.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	_, err = m.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelId, err := strconv.ParseUint(c.PostForm("model_id"), 10, 64)
	if err != nil || modelId == 0 {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	url, err := m.fm.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if url == "" {
		c.JSON(http.StatusUnprocessableEntity, constants.FileURLError)
		return
	}

	ext := filepath.Ext(file.Filename)

	File := entities.File{
		OwnerId:   userID,
		Title:     file.Filename,
		Url:       url,
		Extension: ext,
	}

	File.Prepare()

	validateErr := File.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	savedFile, err := m.fileApp.SaveFile(&File)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	modelFile, err := m.modelApp.SaveModelFile(savedFile, modelId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusCreated, modelFile)
}
