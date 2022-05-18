package worker

import (
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/task"
)

type Worker struct {
	task func() error
}

func NewWorker() *Worker {
	w := new(Worker)

	switch config.Config.Command {
	case task.HELP:
		w.task = task.Help
	case task.LOG:
		w.task = task.Log
	case task.VERSION:
		w.task = task.Version
	case task.ENV:
		w.task = task.Env
	case task.ZIP:
		w.task = task.Zip
	case task.CLEAN:
		w.task = task.Clean
	case task.HELM:
		w.task = task.Helm
	case task.SERVICE:
		w.task = task.Service
	}

	return w
}

func (w *Worker) DoTask() error {
	return w.task()
}
