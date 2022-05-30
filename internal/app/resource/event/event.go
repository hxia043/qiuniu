package event

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var eventUrlPattern string = "%s/api/v1/namespaces/%s/events/"

type Event struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (e *Event) Log() error {
	fmt.Println("Info: collect event log start...")

	collector := collector.NewCollector()

	eventUrl := fmt.Sprintf(eventUrlPattern, e.host, e.namespace)
	resp, err := collector.CollectLog(e.token, eventUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, e.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect event log finished.")

	return nil
}

func NewEvent(host, token, dir string) *Event {
	return &Event{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "event"),
	}
}
