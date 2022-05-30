package endpoint

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var endpointUrlPattern string = "%s/api/v1/namespaces/%s/endpoints"

type Endpoint struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (e *Endpoint) Log() error {
	fmt.Println("Info: collect endpoint log start...")

	collector := collector.NewCollector()

	endpointUrl := fmt.Sprintf(endpointUrlPattern, e.host, e.namespace)
	resp, err := collector.CollectLog(e.token, endpointUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, e.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect endpoint log finished.")

	return nil
}

func NewEndpoint(host, token, dir string) *Endpoint {
	return &Endpoint{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "endpoint"),
	}
}
