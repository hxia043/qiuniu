package models

import (
	"github/hxia043/qiuniu/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	response.SendHttpResponse(ctx, http.StatusCreated, "pong", true)
}
