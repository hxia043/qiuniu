package pv

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type PersistentVolume struct {
	logDir string
}

func (p *PersistentVolume) Log() error {
	collector := collector.NewCollector()

	pvUrl := fmt.Sprintf(config.PvUrlPattern, request.Request.Host, request.Request.Port)
	resp, err := collector.CollectLog(pvUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, p.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewPersistentVolume(dir string) *PersistentVolume {
	logDir := path.Join(dir, "pv")
	os.MkdirAll(logDir, os.ModePerm)
	return &PersistentVolume{
		logDir: logDir,
	}
}
