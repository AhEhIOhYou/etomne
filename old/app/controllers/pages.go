package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Index",
	})
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusOK, "404.html", gin.H{
		"title": "Not found",
	})
}
