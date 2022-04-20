package job

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Job struct {
	Namespace string
	logDir    string
}

func (j *Job) Log() error {
	collector := collector.NewCollector()

	jobUrl := fmt.Sprintf(config.JobUrlPattern, request.Request.Host, request.Request.Port, j.Namespace)
	resp, err := collector.CollectLog(jobUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, j.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewJob(dir string) *Job {
	logDir := path.Join(dir, "job")
	os.MkdirAll(logDir, os.ModePerm)

	return &Job{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
