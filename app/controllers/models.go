package controllers

import (
	"etomne/app/models"
	"etomne/app/server"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type getModelRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func Models(c *gin.Context) {
	modelsList, err := models.GetAllModels(server.Connect())
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, modelsList)
	return
}

func Model(c *gin.Context) {
	var req getModelRequest
	if err := c.ShouldBindUri(&req); err != nil {
		log.Fatal(err)
		return
	}
	m, err := models.GetModelById(req.ID, server.Connect())
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, m)
	return
}
