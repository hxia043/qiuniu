package endpoint

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var endpointUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/endpoints"

type Endpoint struct {
	namespace string
	logDir    string
}

func (e *Endpoint) Log() error {
	collector := collector.NewCollector()

	endpointUrl := fmt.Sprintf(endpointUrlPattern, request.Request.Host, request.Request.Port, e.namespace)
	resp, err := collector.CollectLog(endpointUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, e.logDir); err != nil {
		return err
	}

	return nil
}

func NewEndpoint(dir string) *Endpoint {
	return &Endpoint{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "endpoint"),
	}
}
