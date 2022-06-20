package controllers

import (
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Index",
	})
}

func NotFound(c *gin.Context) {
	c.HTML(200, "404.html", gin.H{
		"title": "Not found",
	})
}
