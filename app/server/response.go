package server

import "github.com/gin-gonic/gin"

type HTTPResponse struct {
	Code    int            `json:"code" example:"200"`
	Message map[string]any `json:"message"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func NewResponse(ctx *gin.Context, status int, msg map[string]any) {
	ctx.Header("Content-Type", "application/json")
	resp := HTTPResponse{
		Code: status,
		Message: map[string]any{
			"response": msg,
		},
	}
	ctx.JSON(status, resp)
}

func NewError(ctx *gin.Context, status int, err string) {
	ctx.Header("Content-Type", "application/json")
	WriteLog(Error, err)
	er := HTTPError{
		Code:    status,
		Message: err,
	}
	ctx.JSON(status, er)
}
