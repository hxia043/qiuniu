package serviceaccount

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type ServiceAccount struct {
	Namespace string
	logDir    string
}

func (s *ServiceAccount) Log() error {
	collector := collector.NewCollector()

	serviceAccountUrl := fmt.Sprintf(config.ServiceaccountUrlPattern, request.Request.Host, request.Request.Port, s.Namespace)
	resp, err := collector.CollectLog(serviceAccountUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, s.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewServiceAccount(dir string) *ServiceAccount {
	logDir := path.Join(dir, "service_account")
	os.MkdirAll(logDir, os.ModePerm)

	return &ServiceAccount{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
