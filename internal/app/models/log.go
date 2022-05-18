package models

import (
	"encoding/json"
	"errors"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/resource"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/k8splugin"
	"github/hxia043/qiuniu/internal/pkg/request"
	"github/hxia043/qiuniu/internal/pkg/response"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

const KIND string = "logs"

var workspace = path.Join(os.Getenv("HOME"), "qiuniu")
var descriptionFile = path.Join(workspace, "qiuniu_description.json")

type LogRequest struct {
	Namespace  string `json:"namespace,omitempty"`
	Host       string `json:"host,omitempty"`
	Port       string `json:"port,omitempty"`
	Token      string `json:"token,omitempty"`
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

type Logs struct {
	Kind  string `json:"Kind,omitempty"`
	Items []Log  `json:"Items,omitempty"`
}

type Log struct {
	CollectTime string `json:"collect_time,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        string `json:"port,omitempty"`
	Name        string `json:"name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
}

func validateRequestParameters(params *LogRequest) error {
	if params.Host == "" && params.Port == "" && params.Token != "" {
		return errors.New("please configure the host and port")
	}

	if params.Kubeconfig != "" {
		if err := k8splugin.ValidateKubeconfig(params.Kubeconfig); err != nil {
			return err
		}
	}

	return nil
}

func createWorkDir(workspace string) (string, string, error) {
	describeId := time.Now().Format(time.RFC3339)
	describeDir := path.Join(workspace, describeId)
	if err := os.MkdirAll(describeDir, os.ModePerm); err != nil {
		return "", "", err
	}

	return describeId, describeDir, nil
}

func updateDescriptionFile(workspace, id string) error {
	log := Log{
		Name:        id,
		Namespace:   config.Config.Namespace,
		Host:        config.Config.Host,
		Port:        config.Config.Port,
		CollectTime: id,
	}

	text, err := json.MarshalIndent(log, "", "    ")
	if err != nil {
		return err
	}

	descriptionInfo := new(Logs)
	f := file.NewFile(descriptionFile, text)
	if !f.Exist() {
		descriptionInfo.Kind = KIND
		descriptionInfo.Items = append(descriptionInfo.Items, log)
	} else {
		descriptionInfoBytes, err := ioutil.ReadFile(f.Path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(descriptionInfoBytes, descriptionInfo); err != nil {
			return err
		}

		descriptionInfo.Items = append(descriptionInfo.Items, log)
	}

	f.Text, err = json.MarshalIndent(descriptionInfo, "", "    ")
	if err != nil {
		return err
	}

	if err := f.WriteFile(); err != nil {
		return err
	}

	return nil
}

func InitConfigFromLogRequest(request *LogRequest) {
	config.Config.Host = request.Host
	config.Config.Port = request.Port
	config.Config.Namespace = request.Namespace
	config.Config.Token = request.Token
}

func CollectLog(ctx *gin.Context) {
	logRequest := new(LogRequest)
	if err := ctx.ShouldBindJSON(logRequest); err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	if err := validateRequestParameters(logRequest); err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	// The purpose is like a parser for config
	InitConfigFromLogRequest(logRequest)
	request.InitLogRequest(config.Config.Host, config.Config.Port, config.Config.Token, false)

	id, workdir, err := createWorkDir(workspace)
	if err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	updateDescriptionFile(workspace, id)

	resources := resource.NewResources(workdir)
	done := make(chan error)
	resourceCount := 0

	for _, res := range resources {
		go func(r resource.Resource) {
			err := r.Log()
			done <- err
		}(res)
	}

	for err := range done {
		if err != nil {
			response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
		}

		resourceCount += 1
		if resourceCount == len(resources) {
			response.SendHttpResponse(ctx, http.StatusCreated, "log collect finished", true)
			break
		}
	}
}
