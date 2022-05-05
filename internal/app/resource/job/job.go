package job

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var jobUrlPattern string = "https://%s:%s/apis/batch/v1/namespaces/%s/jobs"

type Job struct {
	namespace string
	logDir    string
}

func (j *Job) Log() error {
	collector := collector.NewCollector()

	jobUrl := fmt.Sprintf(jobUrlPattern, request.Request.Host, request.Request.Port, j.namespace)
	resp, err := collector.CollectLog(jobUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, j.logDir); err != nil {
		return err
	}

	return nil
}

func NewJob(dir string) *Job {
	return &Job{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "job"),
	}
}
