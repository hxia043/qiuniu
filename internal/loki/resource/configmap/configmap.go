package configmap

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type ConfigMap struct {
	Namespace string
	logDir    string
}

func (c *ConfigMap) Log() error {
	collector := collector.NewCollector()

	configMapUrl := fmt.Sprintf(config.ConfigmapUrlPattern, request.Request.Host, request.Request.Port, c.Namespace)
	resp, err := collector.CollectLog(configMapUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, c.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewConfigMap(dir string) *ConfigMap {
	logDir := path.Join(dir, "configmap")
	os.MkdirAll(logDir, os.ModePerm)

	return &ConfigMap{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
