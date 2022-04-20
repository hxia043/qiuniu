package endpoint

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Endpoint struct {
	Namespace string
	logDir    string
}

func (e *Endpoint) Log() error {
	collector := collector.NewCollector()

	endpointUrl := fmt.Sprintf(config.EndpointUrlPattern, request.Request.Host, request.Request.Port, e.Namespace)
	resp, err := collector.CollectLog(endpointUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, e.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewEndpoint(dir string) *Endpoint {
	logDir := path.Join(dir, "endpoint")
	os.MkdirAll(logDir, os.ModePerm)

	return &Endpoint{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
