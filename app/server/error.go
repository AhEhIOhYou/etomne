package server

import "github.com/gin-gonic/gin"

func NewError(ctx *gin.Context, status int, err string) {
	er := HTTPError{
		Code:    status,
		Message: err,
	}
	ctx.JSON(status, er)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
