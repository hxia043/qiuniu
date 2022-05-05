package event

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var eventUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/events/"

type Event struct {
	namespace string
	logDir    string
}

func (e *Event) Log() error {
	collector := collector.NewCollector()

	eventUrl := fmt.Sprintf(eventUrlPattern, request.Request.Host, request.Request.Port, e.namespace)
	resp, err := collector.CollectLog(eventUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, e.logDir); err != nil {
		return err
	}

	return nil
}

func NewEvent(dir string) *Event {
	return &Event{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "event"),
	}
}
