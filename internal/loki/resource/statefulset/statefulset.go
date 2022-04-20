package statefulset

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Statefulset struct {
	Namespace string
	logDir    string
}

func (s *Statefulset) Log() error {
	collector := collector.NewCollector()
	statefulsetUrl := fmt.Sprintf(config.StatefulsetUrlPattern, request.Request.Host, request.Request.Port, s.Namespace)

	resp, err := collector.CollectLog(statefulsetUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, s.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewStatefulset(dir string) *Statefulset {
	logDir := path.Join(dir, "statefulset")
	os.MkdirAll(logDir, os.ModePerm)

	return &Statefulset{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
