package storageclass

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var storageclassUrlPattern string = "%s/apis/storage.k8s.io/v1/storageclasses"

type StorageClass struct {
	host   string
	token  string
	logDir string
}

func (s *StorageClass) Log() error {
	fmt.Println("Info: collect storageclass log start...")

	collector := collector.NewCollector()

	storageClassUrl := fmt.Sprintf(storageclassUrlPattern, s.host)
	resp, err := collector.CollectLog(s.token, storageClassUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect storageclass log finished.")

	return nil
}

func NewStorageClass(host, token, dir string) *StorageClass {
	return &StorageClass{
		host:   host,
		token:  token,
		logDir: path.Join(dir, "storage_class"),
	}
}
