package statefulset

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var statefulsetUrlPattern string = "%s/apis/apps/v1/namespaces/%s/statefulsets/"

type Statefulset struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (s *Statefulset) Log() error {
	fmt.Println("Info: collect statefulset log start...")

	collector := collector.NewCollector()

	statefulsetUrl := fmt.Sprintf(statefulsetUrlPattern, s.host, s.namespace)
	resp, err := collector.CollectLog(s.token, statefulsetUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect statefulset log finished.")

	return nil
}

func NewStatefulset(host, token, dir string) *Statefulset {
	return &Statefulset{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "statefulset"),
	}
}
