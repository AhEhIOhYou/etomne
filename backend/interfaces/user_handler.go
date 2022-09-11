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

// SaveUser godoc
// @Summary     Save user
// @Tags        Users
// @Param       name		  json      string  true  "Username"
// @Param       email         json      string  true  "User email"
// @Param       password      json      string  true  "User password"
// @Success     201   {object} entities.PublicUser
// @Failure     422   {string} string  "invalid_json"
// @Failure     500   {string} string  "error"
// @Router      /users [post]
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

// GetUsers godoc
// @Summary     Get users
// @Tags        Users
// @Produce		json
// @Param		user_id  query  string  false   "User ID"
// @Param		count    query  string  false   "Count"
// @Success     200  {object} []entities.PublicUser
// @Failure     400  {string} string  "invalid query"
// @Failure     500  {string} string  "error"
// @Router      /users [get]
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

// SaveUserPhoto doc
// @Summary		Save user photo
// @Tags		Users
// @Accept		mpfd
// @Produce		json
// @Param		size  formData  string  false "Size"
// @Param		file  formData  file    true  "File"
// @Success		201  {object}  entities.File
// @Failure     401  string  unauthorized
// @Failure     400  string  error
// @Failure     422  string  error
// @Failure     500  string  error
// @Router		/users/addfile/ [post]
// @Security	bearerAuth
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
