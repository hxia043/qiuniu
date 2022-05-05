package storageclass

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var storageclassUrlPattern string = "https://%s:%s/apis/storage.k8s.io/v1/storageclasses"

type StorageClass struct {
	logDir string
}

func (s *StorageClass) Log() error {
	collector := collector.NewCollector()

	storageClassUrl := fmt.Sprintf(storageclassUrlPattern, request.Request.Host, request.Request.Port)
	resp, err := collector.CollectLog(storageClassUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, s.logDir); err != nil {
		return err
	}

	return nil
}

func NewStorageClass(dir string) *StorageClass {
	return &StorageClass{
		logDir: path.Join(dir, "storage_class"),
	}
}
