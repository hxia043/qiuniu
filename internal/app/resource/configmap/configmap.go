package configmap

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var configmapUrlPattern string = "%s/api/v1/namespaces/%s/configmaps"

type ConfigMap struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (c *ConfigMap) Log() error {
	fmt.Println("Info: collect configmap log start...")

	collector := collector.NewCollector()

	configmapUrl := fmt.Sprintf(configmapUrlPattern, c.host, c.namespace)
	resp, err := collector.CollectLog(c.token, configmapUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, c.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect configmap log finished.")

	return nil
}

func NewConfigMap(host, token, dir string) *ConfigMap {
	return &ConfigMap{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "configmap"),
	}
}
