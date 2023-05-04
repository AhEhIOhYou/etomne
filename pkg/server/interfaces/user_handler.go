package interfaces

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AhEhIOhYou/etomne/pkg/server/application"
	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/pkg/server/interfaces/filemanager"
	"github.com/gin-gonic/gin"
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

//	@Summary	Save user
//	@Tags		user
//	@Produce	json
//	@Param		data	body		entities.UserRequest	true	"User data"
//	@Success	201		{object}	entities.PublicUser
//	@Failure	422		string		string
//	@Failure	500		string		string
//	@Router		/users [post]
func (s *Users) SaveUser(c *gin.Context) {
	var userReq entities.UserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf(constants.Failed, err))
		return
	}

	validateErr := userReq.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	newUser := userReq.NewUser()
	newUser.Prepare()

	savedUser, err := s.userApp.SaveUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusCreated, savedUser.PublicUser())
}

//	@Summary	Get user data by ID
//	@Tags		user
//	@Produce	json
//	@Param		id	path		int	true	"User ID"
//	@Success	200	{object}	entities.PublicUser
//	@Failure	422	string		string
//	@Failure	500	string		string
//	@Router		/users/{id} [get]
func (s *Users) GetUserByID(c *gin.Context) {
	var err error

	var userID uint64
	userIdQuery := c.Param("user_id")
	userID, err = strconv.ParseUint(userIdQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, constants.UserIDInvalid)
		return
	}

	user, err := s.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, user.PublicUser())
}

// TODO edit user, delete and update user permissions
