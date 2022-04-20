package pvc

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type PersistentVolumeClaims struct {
	Namespace string
	logDir    string
}

func (p *PersistentVolumeClaims) Log() error {
	collector := collector.NewCollector()

	pvcUrl := fmt.Sprintf(config.PvcUrlPattern, request.Request.Host, request.Request.Port, p.Namespace)
	resp, err := collector.CollectLog(pvcUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, p.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewPersistentVolumeClaims(dir string) *PersistentVolumeClaims {
	logDir := path.Join(dir, "pvc")
	os.MkdirAll(logDir, os.ModePerm)

	return &PersistentVolumeClaims{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
