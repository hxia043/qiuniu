package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
	Success bool        `json:"success"`
}

func NewHttpResponse(status int, message interface{}, success bool) *HttpResponse {
	return &HttpResponse{
		Status:  status,
		Message: message,
		Success: success,
	}
}

func SendHttpResponse(ctx *gin.Context, status int, message interface{}, success bool) {
	response := NewHttpResponse(status, message, success)
	ctx.JSON(http.StatusOK, response)
}
