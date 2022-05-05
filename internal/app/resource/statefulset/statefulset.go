package statefulset

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var statefulsetUrlPattern string = "https://%s:%s/apis/apps/v1/namespaces/%s/statefulsets/"

type Statefulset struct {
	namespace string
	logDir    string
}

func (s *Statefulset) Log() error {
	collector := collector.NewCollector()

	statefulsetUrl := fmt.Sprintf(statefulsetUrlPattern, request.Request.Host, request.Request.Port, s.namespace)
	resp, err := collector.CollectLog(statefulsetUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	return nil
}

func NewStatefulset(dir string) *Statefulset {
	return &Statefulset{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "statefulset"),
	}
}
