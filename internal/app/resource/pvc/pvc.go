package pvc

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var pvcUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/persistentvolumeclaims"

type PersistentVolumeClaims struct {
	namespace string
	logDir    string
}

func (p *PersistentVolumeClaims) Log() error {
	collector := collector.NewCollector()

	pvcUrl := fmt.Sprintf(pvcUrlPattern, request.Request.Host, request.Request.Port, p.namespace)
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
	return &PersistentVolumeClaims{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "pvc"),
	}
}
