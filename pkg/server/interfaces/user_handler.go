package interfaces

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AhEhIOhYou/etomne/pkg/server/application"
	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/utils"
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
func (us *Users) SaveUser(c *gin.Context) {
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

	savedUser, err := us.userApp.SaveUser(newUser)
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
func (us *Users) GetUserByID(c *gin.Context) {
	userIdQuery := c.Param("user_id")
	userID, err := strconv.ParseUint(userIdQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, constants.UserIDInvalid)
		return
	}

	user, err := us.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, user.PublicUser())
}

//	@Summary	Update user data
//	@Tags		user
//	@Produce	json
//	@Param		user_id	path		int						true	"User ID"
//	@Param		data	body		entities.UserRequest	true	"User updated data"
//	@Success	200		{object}	entities.User
//	@Failure	400		string		string
//	@Failure	401		string		string
//	@Failure	500		string		string
//	@Router		/users/update/{user_id} [post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (us *Users) UpdateUser(c *gin.Context) {
	metadata, err := us.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	currentUserID, err := us.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	currentUser, err := us.userApp.GetUser(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	updatableUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	if !utils.AccessVerification(updatableUserID, currentUser, false) {
		c.JSON(http.StatusUnauthorized, constants.NotEnoughRights)
		return
	}

	updatableUser, err := us.userApp.GetUser(updatableUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	var userReq entities.UserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	if (userReq == entities.UserRequest{}) {
		c.JSON(http.StatusOK, updatableUser)
		return
	}

	if len(userReq.Email) != 0 {
		updatableUser.Email = userReq.Email
	}

	if len(userReq.Name) != 0 {
		updatableUser.Name = userReq.Name
	}

	updatableUser.BeforeUpdate()

	validateErr := updatableUser.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	updatedUser, err := us.userApp.UpdateUser(updatableUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

//	@Summary	Delete user by ID
//	@Tags		user
//	@Param		user_id	path	int	true	"User ID"
//	@Success	200		string	string
//	@Failure	400		string	string
//	@Failure	401		string	string
//	@Failure	500		string	string
//	@Router		/users/{user_id} [delete]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (us *Users) DeleteUser(c *gin.Context) {
	metadata, err := us.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	currentUserID, err := us.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	currentUser, err := us.userApp.GetUser(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	deletableUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	if !utils.AccessVerification(deletableUserID, currentUser, false) {
		c.JSON(http.StatusUnauthorized, constants.NotEnoughRights)
		return
	}

	err = us.userApp.DeleteUser(deletableUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, constants.DeletedSuccessful)
}

//	@Summary	Update user admin rights
//	@Tags		user
//	@Produce	json
//	@Param		user_id	path		int						true	"User ID"
//	@Param		data	body		entities.UserRequest	true	"User updated data"
//	@Success	200		{object}	entities.User
//	@Failure	400		string		string
//	@Failure	401		string		string
//	@Failure	500		string		string
//	@Router		/users/update/admin/{user_id} [post]
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (us *Users) UpdateUserAdminRights(c *gin.Context) {
	metadata, err := us.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	currentUserID, err := us.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	currentUser, err := us.userApp.GetUser(currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	updatableUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	if !utils.AccessVerification(updatableUserID, currentUser, true) {
		c.JSON(http.StatusUnauthorized, constants.NotEnoughRights)
		return
	}

	updatableUser, err := us.userApp.GetUser(updatableUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf(constants.Failed, err))
		return
	}

	var userReq entities.UserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	if (userReq == entities.UserRequest{}) {
		c.JSON(http.StatusOK, updatableUser)
		return
	}

	updatableUser.IsAdmin = userReq.IsAdmin

	validateErr := updatableUser.Validate()
	if len(validateErr) > 0 {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, validateErr))
		return
	}

	updatedUser, err := us.userApp.UpdateUser(updatableUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
