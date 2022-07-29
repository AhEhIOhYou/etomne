package interfaces

import (
	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/filemanager"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type File struct {
	modelApp application.ModelAppInterface
	fileApp  application.FileAppInterface
	userApp  application.UserAppInterface
	fm       filemanager.ManagerFileInterface
	rd       auth.AuthInterface
	tk       auth.TokenInterface
}

func NewFile(mApp application.ModelAppInterface, uApp application.UserAppInterface, fApp application.FileAppInterface, fm filemanager.ManagerFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *File {
	return &File{
		modelApp: mApp,
		userApp:  uApp,
		fileApp:  fApp,
		fm:       fm,
		rd:       rd,
		tk:       tk,
	}
}

// SaveFile doc
// @Summary		Save uploaded file
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
// @Router		/file [post]
// @Security	bearerAuth
func (f *File) SaveFile(c *gin.Context) {
	metadata, err := f.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := f.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = f.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	url, err := f.fm.UploadFile(file)
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
		OwnerId: userId,
		Title:   file.Filename,
		Url:     url,
	}

	File.Prepare()

	saveFile, saveModelErr := f.fileApp.SaveFile(&File)
	if len(saveModelErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveModelErr)
		return
	}

	var modelId uint64
	modelIdForm := c.PostForm("model_id")
	if modelIdForm != "" {
		modelId, err = strconv.ParseUint(modelIdForm, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid query")
			return
		}
	}

	if modelId != 0 {
		//Добавление файла и модельки в бд model_files
		ModelFile := entities.ModelFile{
			ModelId: modelId,
			FileId:  saveFile.ID,
		}

		_, saveModelErr = f.modelApp.AddModelFile(&ModelFile)
		if len(saveModelErr) > 0 {
			c.JSON(http.StatusUnprocessableEntity, saveModelErr)
			return
		}
	}

	c.JSON(http.StatusCreated, saveFile)
}

// RemoveFile doc
// @Summary		Remove uploaded file
// @Tags		File
// @Accept		mpfd
// @Produce		json
// @Param		file_id   path    string  true  "File ID"
// @Param		model_id  query   string  false "Model ID"
// @Success		201  {object}  entities.File
// @Failure     401  string  unauthorized
// @Failure     400  string  error
// @Failure     422  string  error
// @Failure     500  string  error
// @Router		/file/{id} [post]
// @Security	bearerAuth
func (f *File) RemoveFile(c *gin.Context) {
	metadata, err := f.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := f.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = f.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	fileId, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	var modelId uint64
	modelIdForm := c.Query("model_id")
	if modelIdForm != "" {
		modelId, err = strconv.ParseUint(modelIdForm, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid query")
			return
		}
	}

	isAvaliable, err := f.fileApp.CheckAvailabilityFile(fileId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !isAvaliable {
		c.JSON(http.StatusInternalServerError, "file is unavailable")
		return
	}

	file, err := f.fileApp.GetFile(fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if modelId != 0 {
		//Отвязать файл от модели
		err = f.modelApp.DeleteModelFile(file.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	//Удалить файл с бд
	err = f.fileApp.DeleteFile(file.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//Удалить файл с хранилища
	err = f.fm.DeleteFile(file.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, file)
}
