package models

import (
	"encoding/json"
	"github/hxia043/qiuniu/internal/app/resource"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/k8splugin"
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
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

type Logs struct {
	Kind  string `json:"Kind,omitempty"`
	Items []Log  `json:"Items,omitempty"`
}

type Log struct {
	CollectTime string `json:"collect_time,omitempty"`
	Host        string `json:"host,omitempty"`
	Name        string `json:"name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
}

func createWorkDir() (string, string, error) {
	descriptionId := time.Now().Format(time.RFC3339)

	descriptionDir := path.Join(os.Getenv("HOME"), "qiuniu", descriptionId)
	if err := os.MkdirAll(descriptionDir, os.ModePerm); err != nil {
		return "", "", err
	}

	return descriptionId, descriptionDir, nil
}

func updateDescriptionFile(workspace, id, host, namespace string) error {
	log := Log{
		Name:        id,
		Namespace:   namespace,
		Host:        host,
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

func transferKubeconfig(ctx *gin.Context) ([]byte, error) {
	kubeconfigHeader, err := ctx.FormFile("kubeconfig")
	if err != nil {
		return nil, err
	}

	kubeconfigFile, err := kubeconfigHeader.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		kubeconfigFile.Close()
	}()

	kubeconfig := make([]byte, kubeconfigHeader.Size)
	if _, err := kubeconfigFile.Read(kubeconfig); err != nil {
		return nil, err
	}

	if err := k8splugin.ValidateKubeconfig(kubeconfig); err != nil {
		return nil, err
	}

	return kubeconfig, nil
}

func CollectLog(ctx *gin.Context) {
	logRequest := new(LogRequest)

	form, err := ctx.MultipartForm()
	if err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	if namespace, ok := form.Value["namespace"]; ok && len(namespace) > 0 {
		logRequest.Namespace = namespace[0]
	}

	kubeconfig, err := transferKubeconfig(ctx)
	if err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	namespace, host, token, err := k8splugin.ParseKubeconfig(kubeconfig)
	if err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	if logRequest.Namespace != "" {
		logRequest.Namespace = namespace
	}

	id, workdir, err := createWorkDir()
	if err != nil {
		response.SendHttpResponse(ctx, http.StatusForbidden, err.Error(), false)
	}

	updateDescriptionFile(workspace, id, host, logRequest.Namespace)

	resources := resource.NewResources(host, token, workdir, kubeconfig)
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
