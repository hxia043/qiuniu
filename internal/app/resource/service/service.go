package service

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var serviceUrlPattern string = "%s/api/v1/namespaces/%s/services"

type Service struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (s *Service) Log() error {
	fmt.Println("Info: collect service log start...")

	collector := collector.NewCollector()

	serviceUrl := fmt.Sprintf(serviceUrlPattern, s.host, s.namespace)
	resp, err := collector.CollectLog(s.token, serviceUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect service log finished.")

	return nil
}

func NewService(host, token, dir string) *Service {
	return &Service{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "service"),
	}
}
