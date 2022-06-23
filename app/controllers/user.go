package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"etomne/app/entities"
	"etomne/app/models"
	"etomne/app/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func User(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "User Page",
	})
}

func UserLogin(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "User Page",
		})
	} else {
		login := c.PostForm("login")
		pass := c.PostForm("pass")
		hash := md5.Sum([]byte(pass))
		hashedPass := hex.EncodeToString(hash[:])

		userModel := models.UserModel{
			Db: server.Connect(),
		}

		user, _ := userModel.Login(login, hashedPass)

		if user == (entities.User{}) {
			c.JSON(http.StatusOK, gin.H{
				"No": "no",
			})
		} else {

			timeInt := strconv.FormatInt(time.Now().Unix(), 10)
			token := login + pass + timeInt
			hashToken := md5.Sum([]byte(token))
			hashedToken := hex.EncodeToString(hashToken[:])
			livingTime := int(60 * time.Minute)

			c.SetCookie("token", hashedToken, livingTime, "/", "localhost", false, true)

			c.Redirect(http.StatusPermanentRedirect, "/")
		}
	}
}
func UserReg(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "reg.html", gin.H{
			"title": "Registration",
		})
	} else {
		login := c.PostForm("login")
		pass := c.PostForm("pass")
		name := c.PostForm("name")
		hash := md5.Sum([]byte(pass))
		hashedPass := hex.EncodeToString(hash[:])

		userModel := models.UserModel{
			Db: server.Connect(),
		}

		id, _ := userModel.Create(login, hashedPass, name)

		if id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"No": "no",
			})
		} else {
			c.Redirect(http.StatusPermanentRedirect, "/login")
		}
	}
}
func UserLogout(c *gin.Context) {
	c.SetCookie("token", " ", 1, "/", "localhost", false, true)
	c.Redirect(http.StatusPermanentRedirect, "/")
}
