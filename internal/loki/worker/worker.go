package worker

import (
	"fmt"
	"log"
	"loki/internal/loki/config"
	rconfig "loki/internal/loki/resource/config"
	"loki/internal/loki/task"
	"loki/internal/loki/version"
	"loki/internal/pkg/clean"
	"loki/internal/pkg/request"
	"loki/internal/pkg/zip"
	"os"
	"path"
	"runtime"
)

type Worker struct {
	task func() error
}

func NewWorker(cfg config.Config) *Worker {
	w := new(Worker)

	switch cfg.Command {
	case task.HELP:
		w.task = task.Help
	case task.LOG:
		parseAppLogConfig(cfg)
		w.task = task.Log
	case task.VERSION:
		parseAppVersionConfig(cfg)
		w.task = task.Version
	case task.ENV:
		w.task = task.Env
	case task.ZIP:
		parseAppZipDirConfig(cfg.ZipDir)
		w.task = task.Zip
	case task.CLEAN:
		parseAppCleanConfig(cfg)
		w.task = task.Clean
	case task.HELM:
		parseAppHelmConfig(cfg)
		w.task = task.Helm
	}

	return w
}

func parseAppHelmConfig(cfg config.Config) {
	rconfig.ResourceConfig.KubeConfig = cfg.KubeConfig
	rconfig.ResourceConfig.Namespace = cfg.Namespace

	parseAppWorkspaceConfig(cfg.WorkSpace)
}

func parseAppCleanConfig(cfg config.Config) {
	if cfg.Interval != 0 {
		clean.Interval = cfg.Interval
	}

	if cfg.WorkSpace != "" {
		clean.Workdir = path.Join(cfg.WorkSpace, "loki")
	} else {
		err := fmt.Errorf("loki: no workspace defined in command line or env :(")
		log.Fatal(err)
	}
}

func parseAppVersionConfig(cfg config.Config) {
	version.Version.LokiVersion = cfg.Version
	version.Version.LokiType = cfg.Type
	version.Version.GoVersion = runtime.Version()
}

func parseAppLogConfig(cfg config.Config) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/yaml"
	headers["Authorization"] = "Bearer " + cfg.Token

	request.Request.Host = cfg.Host
	request.Request.Port = cfg.Port
	request.Request.IsVerify = cfg.IsVerify
	request.Request.Headers = headers
	request.Request.Method = string(request.GET_REQUEST)

	rconfig.ResourceConfig.Namespace = cfg.Namespace

	parseAppWorkspaceConfig(cfg.WorkSpace)
}

func parseAppWorkspaceConfig(workspace string) {
	if workspace != "" {
		rconfig.ResourceConfig.Workdir = path.Join(workspace, "loki")
	} else {
		defaultWorkspace := os.Getenv("HOME")
		rconfig.ResourceConfig.Workdir = path.Join(defaultWorkspace, "loki")
	}
}

func parseAppZipDirConfig(dir string) {
	if dir != "" {
		zip.Dir = dir
	}
}

func (w *Worker) DoTask() error {
	err := w.task()
	if err != nil {
		return err
	}

	return nil
}
