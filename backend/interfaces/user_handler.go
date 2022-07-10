package interfaces

import (
	"etomne/backend/application"
	"etomne/backend/domain/entities"
	auth2 "etomne/backend/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Users struct {
	us application.UserAppInterface
	rd auth2.AuthInterface
	tk auth2.TokenInterface
}

func NewUsers(us application.UserAppInterface, rd auth2.AuthInterface, tk auth2.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (s *Users) SaveUser(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid_json",
		})
		return
	}

	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, validateErr)
		return
	}

	newUser, err := s.us.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, newUser.PublicUser())
}

// GetUsers godoc
// @Summary      Get users
// @Description  Get users
// @Tags         Users
// @Success      200  {object}  entities.PublicUser
// @Router       /users [get]
func (s *Users) GetUsers(c *gin.Context) {
	users := entities.Users{}
	var err error
	users, err = s.us.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (s *Users) GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.us.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}
