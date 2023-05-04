package interfaces

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/constants"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Authenticate struct {
	us application.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewAuthenticate(uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Authenticate {
	return &Authenticate{
		us: uApp,
		rd: rd,
		tk: tk,
	}
}

//	@Summary	Login user
//	@Tags		Auth
//	@Produce	json
//	@Param		data	body		entities.LoginRequest	true	"User login data"
//	@Success	200		{object}	entities.UserResponse
//	@Failure	422		string		string
//	@Failure	500		string		string
//	@Router		/users/login [post]
func (au *Authenticate) Login(c *gin.Context) {
	var userLogin *entities.LoginRequest

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf(constants.Failed, err))
		return
	}

	userLogin.Prepare()

	user, err := au.us.GetUserByEmailAndPassword(&entities.User{
		Email:    userLogin.Email,
		Password: userLogin.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	ts, err := au.tk.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	err = au.rd.CreateAuth(user.ID, ts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	c.JSON(http.StatusOK, &entities.UserResponse{
		PublicUser: entities.PublicUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		UserAuth: entities.UserAuth{
			RefreshToken: ts.RefreshToken,
			AccessToken:  ts.AccessToken,
		},
	})
}

//	@Summary	Logout user
//	@Tags		Auth
//	@Success	200	{string}	string
//	@Failure	401	string		string
//	@Failure	500	string		string
//	@Router		/users/logout [get]
//	Security	BearerAuth
//	@Param		Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func (au *Authenticate) Logout(c *gin.Context) {
	metadata, err := au.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}
	err = au.rd.DeleteTokens(metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}
	c.JSON(http.StatusOK, constants.LogoutSuccessful)
}

//	@Summary	Refresh user session
//	@Tags		Auth
//	@Produce	json
//	@Param		data	body		entities.UserAuth	true	"User tokens data"
//	@Success	200		{object}	entities.UserResponse
//	@Failure	401		string		string
//	@Failure	403		string		string
//	@Failure	422		string		string
//	@Failure	500		string		string
//	@Router		/users/refresh [post]
func (au *Authenticate) Refresh(c *gin.Context) {
	var tokens *entities.UserAuth

	if err := c.ShouldBindJSON(&tokens); err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf(constants.Failed, err))
		return
	}

	token, err := jwt.Parse(tokens.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf(constants.Failed, err))
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, err))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf(constants.Failed, constants.CannotGetUUID))
			return
		}

		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf(constants.Failed, err))
			return
		}

		err = au.rd.DeleteRefresh(refreshUuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, fmt.Sprintf(constants.Failed, constants.Unauthorized))
			return
		}

		ts, err := au.tk.CreateToken(userID)
		if err != nil {
			c.JSON(http.StatusForbidden, fmt.Sprintf(constants.Failed, err))
			return
		}

		err = au.rd.CreateAuth(userID, ts)
		if err != nil {
			c.JSON(http.StatusForbidden, fmt.Sprintf(constants.Failed, err))
			return
		}

		c.JSON(http.StatusCreated, &entities.UserResponse{
			PublicUser: entities.PublicUser{
				ID: userID,
			},
			UserAuth: entities.UserAuth{
				RefreshToken: ts.RefreshToken,
				AccessToken:  ts.AccessToken,
			},
		})
	} else {
		c.JSON(http.StatusUnauthorized, constants.RefreshTokenExpired)
	}
}
