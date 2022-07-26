package interfaces

import (
	"fmt"
	"github.com/AhEhIOhYou/etomne/backend/application"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
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

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RToken struct {
	RefreshToken string `json:"refresh_token"`
}

type Tokens struct {
}

// Login doc
// @Summary		Login user
// @Tags		Authorization
// @Accept 		json
// @Produce		json
// @Param 		data 	body 	  AuthUser true "Auth user data"
// @Success 	200 	{object}  UserData
// @Failure     422  {string} string  "invalid json provided"
// @Failure     500  {string} string  "error"
// @Router		/users/login [post]
func (au *Authenticate) Login(c *gin.Context) {

	var user *entities.User
	var tokenErr = map[string]string{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}

	u, userErr := au.us.GetUserByEmailAndPassword(user)
	if userErr != nil {
		c.JSON(http.StatusInternalServerError, userErr)
		return
	}

	ts, tErr := au.tk.CreateToken(u.ID)
	if tErr != nil {
		tokenErr["token_error"] = tErr.Error()
		c.JSON(http.StatusInternalServerError, tErr.Error())
		return
	}

	saveErr := au.rd.CreateAuth(u.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr.Error())
		return
	}

	var userData = UserData{
		ID:           u.ID,
		Name:         u.Name,
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	}

	c.JSON(http.StatusOK, userData)
}

// Logout doc
// @Summary		Logout user
// @Tags		Authorization
// @Accept 		json
// @Produce		json
// @Success 	200  {string} string  "successfully logged out"
// @Failure     401  {string} string  "unauthorized"
// @Failure     500  {string} string  "error"
// @Router		/users/logout [post]
func (au *Authenticate) Logout(c *gin.Context) {
	metadata, err := au.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleteErr := au.rd.DeleteTokens(metadata)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, deleteErr.Error())
		return
	}
	c.JSON(http.StatusOK, "successfully logged out")
}

// Refresh doc
// @Summary		Refresh user token
// @Tags		Authorization
// @Accept 		json
// @Produce		json
// @Param 		data  body 	  RToken true "Refresh token"
// @Success 	200  {object} UserData
// @Failure     401  {string} string  "unauthorized"
// @Failure     403  {string} string  "error"
// @Failure     422  {string} string  "error"
// @Failure     500  {string} string  "error"
// @Router		/users/refresh [post]
func (au *Authenticate) Refresh(c *gin.Context) {
	var rToken *RToken
	if err := c.ShouldBindJSON(&rToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	token, err := jwt.Parse(rToken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "cannot get uuid")
			return
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "error occurred")
			return
		}

		delErr := au.rd.DeleteRefresh(refreshUuid)
		if delErr != nil { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		ts, createErr := au.tk.CreateToken(userId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}

		saveErr := au.rd.CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}

		var userData = UserData{
			ID:           userId,
			AccessToken:  ts.AccessToken,
			RefreshToken: ts.RefreshToken,
		}

		c.JSON(http.StatusCreated, userData)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh token expired")
	}
}
