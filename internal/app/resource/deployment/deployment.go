package deployment

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var deploymentUrlPattern string = "https://%s:%s/apis/apps/v1/namespaces/%s/deployments/"

type Deployment struct {
	namespace string
	logDir    string
}

func (d *Deployment) Log() error {
	collector := collector.NewCollector()

	deploymentUrl := fmt.Sprintf(deploymentUrlPattern, request.Request.Host, request.Request.Port, d.namespace)
	resp, err := collector.CollectLog(deploymentUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, d.logDir); err != nil {
		return err
	}

	return nil
}

func NewDeployment(dir string) *Deployment {
	return &Deployment{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "deployment"),
	}
}
