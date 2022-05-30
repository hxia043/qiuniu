package job

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var jobUrlPattern string = "%s/apis/batch/v1/namespaces/%s/jobs"

type Job struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (j *Job) Log() error {
	fmt.Println("Info: collect job log start...")

	collector := collector.NewCollector()

	jobUrl := fmt.Sprintf(jobUrlPattern, j.host, j.namespace)
	resp, err := collector.CollectLog(j.token, jobUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, j.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect job log finished.")

	return nil
}

func NewJob(host, token, dir string) *Job {
	return &Job{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "job"),
	}
}
