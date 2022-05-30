package models

import (
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/clean"
	"github/hxia043/qiuniu/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CleanLog(ctx *gin.Context) {
	if err := clean.Clean(workspace, config.Config.Interval); err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	response.SendHttpResponse(ctx, http.StatusCreated, "log cleaned", true)
}
