package event

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Event struct {
	Namespace string
	logDir    string
}

func (e *Event) Log() error {
	collector := collector.NewCollector()
	eventUrl := fmt.Sprintf(config.EventUrlPattern, request.Request.Host, request.Request.Port, e.Namespace)

	resp, err := collector.CollectLog(eventUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, e.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewEvent(dir string) *Event {
	logDir := path.Join(dir, "event")
	os.MkdirAll(logDir, os.ModePerm)

	return &Event{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
