package interfaces

import (
	"github.com/AhEhIOhYou/etomne/backend/application"
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

type modelId struct {
	ModelId int `json:"model_id"`
}

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

	fileId, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	var mId modelId
	if err := c.ShouldBindJSON(&mId); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid_json",
		})
		return
	}

	isAvaliable, err := f.modelApp.CheckAvailability(fileId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !isAvaliable {
		c.JSON(http.StatusInternalServerError, "file is unavailable")
		return
	}

	/*
		file, err := f.fileApp.GetFile(fileId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
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

		//Отвязать файл от модели
		err = f.modelApp.DeleteModelFile(file.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	*/

	c.JSON(http.StatusOK, gin.H{
		"file id":  fileId,
		"model_id": mId.ModelId,
	})
}

func (f *File) RemoveFile(c *gin.Context) {

}
