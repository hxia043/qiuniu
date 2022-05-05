package configmap

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var configmapUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/configmaps"

type ConfigMap struct {
	namespace string
	logDir    string
}

func (c *ConfigMap) Log() error {
	collector := collector.NewCollector()

	configmapUrl := fmt.Sprintf(configmapUrlPattern, request.Request.Host, request.Request.Port, c.namespace)
	resp, err := collector.CollectLog(configmapUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, c.logDir); err != nil {
		return err
	}

	return nil
}

func NewConfigMap(dir string) *ConfigMap {
	return &ConfigMap{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "configmap"),
	}
}
