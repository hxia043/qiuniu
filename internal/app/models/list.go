package models

import (
	"encoding/json"
	"github/hxia043/qiuniu/internal/pkg/response"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListLog(ctx *gin.Context) {
	file, err := ioutil.ReadFile(descriptionFile)
	if err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	logs := new(Logs)
	if err := json.Unmarshal(file, logs); err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	response.SendHttpResponse(ctx, http.StatusCreated, logs, true)
}
