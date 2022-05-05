package pv

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var pvUrlPattern string = "https://%s:%s/api/v1/persistentvolumes/"

type PersistentVolume struct {
	logDir string
}

func (p *PersistentVolume) Log() error {
	collector := collector.NewCollector()

	pvUrl := fmt.Sprintf(pvUrlPattern, request.Request.Host, request.Request.Port)
	resp, err := collector.CollectLog(pvUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, p.logDir); err != nil {
		return err
	}

	return nil
}

func NewPersistentVolume(dir string) *PersistentVolume {
	return &PersistentVolume{
		logDir: path.Join(dir, "pv"),
	}
}
