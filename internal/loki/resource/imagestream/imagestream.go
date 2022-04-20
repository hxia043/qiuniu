package imagestream

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type ImageStream struct {
	Namespace string
	logDir    string
}

func (i *ImageStream) Log() error {
	collector := collector.NewCollector()

	imageStreamUrl := fmt.Sprintf(config.ImagestreamUrlPattern, request.Request.Host, request.Request.Port, i.Namespace)
	resp, err := collector.CollectLog(imageStreamUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, i.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewImageStream(dir string) *ImageStream {
	logDir := path.Join(dir, "image_stream")
	os.MkdirAll(logDir, os.ModePerm)

	return &ImageStream{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
