package imagestream

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var imagestreamUrlPattern string = "%s/apis/image.openshift.io/v1/namespaces/%s/imagestreams"

type ImageStream struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (i *ImageStream) Log() error {
	fmt.Println("Info: collect image log start...")

	collector := collector.NewCollector()

	imageStreamUrl := fmt.Sprintf(imagestreamUrlPattern, i.host, i.namespace)
	resp, err := collector.CollectLog(i.token, imageStreamUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, i.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect image log finished.")

	return nil
}

func NewImageStream(host, token, dir string) *ImageStream {
	return &ImageStream{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "image_stream"),
	}
}
