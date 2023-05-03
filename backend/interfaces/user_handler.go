package interfaces

import (
	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/filemanager"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

type Users struct {
	userApp application.UserAppInterface
	fm      filemanager.ManagerFileInterface
	rd      auth.AuthInterface
	tk      auth.TokenInterface
}

func NewUsers(userApp application.UserAppInterface, fm filemanager.ManagerFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		userApp: userApp,
		fm:      fm,
		rd:      rd,
		tk:      tk,
	}
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary      Save user
// @Tags         user
// @Produce      json
// @Param        data body entities.User true "user data"
// @Success      201  {object}   entities.PublicUser
// @Failure      422  string    string
// @Failure      500  string    string
// @Router       /users [post]
func (s *Users) SaveUser(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid_json")
		return
	}

	user.Prepared()

	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, validateErr)
		return
	}

	newUser, err := s.userApp.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (s *Users) GetUsers(c *gin.Context) {
	users := entities.Users{}
	var err error
	var ok bool

	var userId uint64
	userIdQuery, ok := c.GetQuery("user_id")
	if ok {
		userId, err = strconv.ParseUint(userIdQuery, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid query")
			return
		}
	}

	var count uint64
	countQuery, ok := c.GetQuery("count")
	if ok {
		count, err = strconv.ParseUint(countQuery, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid query")
			return
		}
	}

	if count == 0 {
		count = 2
	}

	if userId != 0 {
		user, err := s.userApp.GetUser(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		users = append(users, *user)
	} else {
		users, err = s.userApp.GetUsers(count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (s *Users) SaveUserPhoto(c *gin.Context) {
	metadata, err := s.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, err := s.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = s.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	photoSize, err := strconv.ParseUint(c.PostForm("photo_size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid query")
		return
	}
	if photoSize == 0 {
		photoSize = 200
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	url, err := s.fm.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if url == "" {
		c.JSON(http.StatusUnprocessableEntity, "something went wrong")
		return
	}

	ext := filepath.Ext(file.Filename)

	File := entities.File{
		OwnerId:   userId,
		Title:     file.Filename,
		Url:       url,
		Extension: ext,
	}

	File.Prepare()
	saveFileErr := File.Validate("")
	if len(saveFileErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveFileErr)
		return
	}

	userPhoto, saveFileErr := s.userApp.SaveUserPhoto(&File, userId, 100)
	if len(saveFileErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveFileErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file":       File,
		"user_photo": userPhoto,
	})
}
