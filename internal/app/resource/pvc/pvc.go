package pvc

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var pvcUrlPattern string = "%s/api/v1/namespaces/%s/persistentvolumeclaims"

type PersistentVolumeClaims struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (p *PersistentVolumeClaims) Log() error {
	fmt.Println("Info: collect persistent volume claim log start...")

	collector := collector.NewCollector()

	pvcUrl := fmt.Sprintf(pvcUrlPattern, p.host, p.namespace)
	resp, err := collector.CollectLog(p.token, pvcUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, p.logDir)
	if err != nil {
		return err
	}

	fmt.Println("Info: collect persistent volume claim log finished.")

	return nil
}

func NewPersistentVolumeClaims(host, token, dir string) *PersistentVolumeClaims {
	return &PersistentVolumeClaims{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "pvc"),
	}
}
