package interfaces

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/constants"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/filemanager"
	"github.com/gin-gonic/gin"
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

//	@Summary	Save file
//	@Tags		file
//	@Produce	json
//	@Param		file	formData	file	true	"Body with files"
//	@Success	201		{object}	entities.File
//	@Failure	400		string		string
//	@Failure	401		string		string
//	@Failure	500		string		string
//	@Router		/file [post]
//	Security	BearerAuth
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (f *File) SaveFile(c *gin.Context) {
	metadata, err := f.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := f.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	_, err = f.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	url, err := f.fm.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}
	if url == "" {
		c.JSON(http.StatusInternalServerError, constants.FileURLError)
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

	savedFile, err := f.fileApp.SaveFile(&File)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusCreated, savedFile)
}

//	@Summary	Delete file by ID
//	@Tags		file
//	@Param		file_id	path		int	true	"File ID"
//	@Success	200		{object}	entities.File
//	@Failure	400		string		string
//	@Failure	401		string		string
//	@Failure	500		string		string
//	@Router		/file/{file_id} [delete]
//	Security	BearerAuth
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (f *File) RemoveFile(c *gin.Context) {
	metadata, err := f.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	userID, err := f.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	_, err = f.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	isAvaliable, err := f.fileApp.CheckAvailabilityFile(fileID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}
	if !isAvaliable {
		c.JSON(http.StatusInternalServerError, constants.FileNotAvaliable)
		return
	}

	file, err := f.fileApp.GetFile(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	// Delete from DB
	err = f.fileApp.DeleteFile(file.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	// Delete from storage
	err = f.fm.DeleteFile(file.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, file)
}
