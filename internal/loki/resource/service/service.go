package service

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Service struct {
	Namespace string
	logDir    string
}

func (s *Service) Log() error {
	collector := collector.NewCollector()

	serviceUrl := fmt.Sprintf(config.ServiceUrlPattern, request.Request.Host, request.Request.Port, s.Namespace)
	resp, err := collector.CollectLog(serviceUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, s.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewService(dir string) *Service {
	logDir := path.Join(dir, "service")
	os.MkdirAll(logDir, os.ModePerm)

	return &Service{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
