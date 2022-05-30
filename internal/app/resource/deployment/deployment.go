package deployment

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var deploymentUrlPattern string = "%s/apis/apps/v1/namespaces/%s/deployments/"

type Deployment struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (d *Deployment) Log() error {
	fmt.Println("Info: collect deployment log start...")

	collector := collector.NewCollector()

	deploymentUrl := fmt.Sprintf(deploymentUrlPattern, d.host, d.namespace)
	resp, err := collector.CollectLog(d.token, deploymentUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, d.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect deployment log finished.")

	return nil
}

func NewDeployment(host, token, dir string) *Deployment {
	return &Deployment{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "deployment"),
	}
}
