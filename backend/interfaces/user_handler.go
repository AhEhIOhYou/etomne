package interfaces

import (
	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Users struct {
	us application.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewUsers(us application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
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
		user, err := s.us.GetUser(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		users = append(users, *user)
	} else {
		users, err = s.us.GetUsers(count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}
