package models

import (
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/response"
	"github/hxia043/qiuniu/internal/pkg/zip"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func TransferLogToData(path string) ([]byte, error) {
	fileName, err := zip.Zip(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := os.Remove(fileName); err != nil {
		return nil, err
	}

	return data, nil
}

func DownloadLog(ctx *gin.Context) {
	name := ctx.Query("name")
	if name != "" {
		f := file.NewFile(descriptionFile, nil)
		if !f.Exist() {
			response.SendHttpResponse(ctx, http.StatusForbidden, "log not found", false)
		}

		data, err := ioutil.ReadFile(f.Path)
		if err != nil {
			response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
		}

		logs := new(Logs)
		err = json.Unmarshal(data, logs)
		if err != nil {
			response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
		}

		for _, item := range logs.Items {
			if item.Name == name {
				filePath := path.Join(workspace, item.Name)
				data, err := TransferLogToData(filePath)
				if err != nil {
					response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
				}

				fileContentDisposition := fmt.Sprintf("attachment; filename=\"%s.zip\"", filePath)
				ctx.Header("Content-Type", "application/octet-stream")
				ctx.Header("Content-Disposition", fileContentDisposition)
				ctx.Data(http.StatusCreated, "application/octet-stream", data)

				break
			}
		}
	} else {
		f := file.NewFile(descriptionFile, nil)
		if !f.Exist() {
			response.SendHttpResponse(ctx, http.StatusForbidden, "log not found", false)
		}

		data, err := TransferLogToData(workspace)
		if err != nil {
			response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
		}

		fileContentDisposition := fmt.Sprintf("attachment; filename=\"%s.zip\"", workspace)
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", fileContentDisposition)
		ctx.Data(http.StatusCreated, "application/octet-stream", data)
	}
}
