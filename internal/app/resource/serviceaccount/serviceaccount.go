package serviceaccount

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var serviceaccountUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/serviceaccounts"

type ServiceAccount struct {
	namespace string
	logDir    string
}

func (s *ServiceAccount) Log() error {
	collector := collector.NewCollector()

	serviceAccountUrl := fmt.Sprintf(serviceaccountUrlPattern, request.Request.Host, request.Request.Port, s.namespace)
	resp, err := collector.CollectLog(serviceAccountUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	return nil
}

func NewServiceAccount(dir string) *ServiceAccount {
	return &ServiceAccount{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "service_account"),
	}
}
