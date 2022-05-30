package pv

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var pvUrlPattern string = "%s/api/v1/persistentvolumes/"

type PersistentVolume struct {
	host   string
	token  string
	logDir string
}

func (p *PersistentVolume) Log() error {
	fmt.Println("Info: collect persistent volume log start...")

	collector := collector.NewCollector()

	pvUrl := fmt.Sprintf(pvUrlPattern, p.host)
	resp, err := collector.CollectLog(p.token, pvUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, p.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect persistent volume log finished.")

	return nil
}

func NewPersistentVolume(host, token, dir string) *PersistentVolume {
	return &PersistentVolume{
		host:   host,
		token:  token,
		logDir: path.Join(dir, "pv"),
	}
}
