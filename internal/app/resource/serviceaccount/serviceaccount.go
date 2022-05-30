package serviceaccount

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var serviceaccountUrlPattern string = "%s/api/v1/namespaces/%s/serviceaccounts"

type ServiceAccount struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (s *ServiceAccount) Log() error {
	fmt.Println("Info: collect service account log start...")

	collector := collector.NewCollector()

	serviceAccountUrl := fmt.Sprintf(serviceaccountUrlPattern, s.host, s.namespace)
	resp, err := collector.CollectLog(s.token, serviceAccountUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect service account log finished.")

	return nil
}

func NewServiceAccount(host, token, dir string) *ServiceAccount {
	return &ServiceAccount{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "service_account"),
	}
}
