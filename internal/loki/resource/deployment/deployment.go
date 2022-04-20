package deployment

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Deployment struct {
	Namespace string
	logDir    string
}

func (d *Deployment) Log() error {
	collector := collector.NewCollector()
	deploymentUrl := fmt.Sprintf(config.DeploymentUrlPattern, request.Request.Host, request.Request.Port, d.Namespace)

	resp, err := collector.CollectLog(deploymentUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, d.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewDeployment(dir string) *Deployment {
	logDir := path.Join(dir, "deployment")
	os.MkdirAll(logDir, os.ModePerm)

	return &Deployment{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
