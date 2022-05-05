package service

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var serviceUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/services"

type Service struct {
	namespace string
	logDir    string
}

func (s *Service) Log() error {
	collector := collector.NewCollector()

	serviceUrl := fmt.Sprintf(serviceUrlPattern, request.Request.Host, request.Request.Port, s.namespace)
	resp, err := collector.CollectLog(serviceUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	return nil
}

func NewService(dir string) *Service {
	return &Service{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "service"),
	}
}
