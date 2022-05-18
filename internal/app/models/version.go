package models

import (
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/response"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type version struct {
	GoVersion     string `json:"go_version"`
	QiuniuType    string `json:"qiuniu_type"`
	QiuniuVersion string `json:"qiuniu_version"`
}

func Version(ctx *gin.Context) {
	v := &version{
		GoVersion:     runtime.Version(),
		QiuniuType:    config.Type,
		QiuniuVersion: config.Version,
	}

	response.SendHttpResponse(ctx, http.StatusCreated, v, true)
}
