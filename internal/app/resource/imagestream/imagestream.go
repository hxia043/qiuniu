package imagestream

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var imagestreamUrlPattern string = "https://%s:%s/apis/image.openshift.io/v1/namespaces/%s/imagestreams"

type ImageStream struct {
	namespace string
	logDir    string
}

func (i *ImageStream) Log() error {
	collector := collector.NewCollector()

	imageStreamUrl := fmt.Sprintf(imagestreamUrlPattern, request.Request.Host, request.Request.Port, i.namespace)
	resp, err := collector.CollectLog(imageStreamUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, i.logDir); err != nil {
		return err
	}

	return nil
}

func NewImageStream(dir string) *ImageStream {
	return &ImageStream{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "image_stream"),
	}
}
