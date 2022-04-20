package storageclass

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type StorageClass struct {
	logDir string
}

func (s *StorageClass) Log() error {
	collector := collector.NewCollector()

	storageClassUrl := fmt.Sprintf(config.StorageclassUrlPattern, request.Request.Host, request.Request.Port)
	resp, err := collector.CollectLog(storageClassUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, s.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewStorageClass(dir string) *StorageClass {
	logDir := path.Join(dir, "storage_class")
	os.MkdirAll(logDir, os.ModePerm)

	return &StorageClass{
		logDir: logDir,
	}
}
