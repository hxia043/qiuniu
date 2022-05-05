package app

import (
	"github/hxia043/qiuniu/internal/app/parser"
	"github/hxia043/qiuniu/internal/app/worker"
)

type App struct{}

func NewApp() (*App, error) {
	p := parser.NewParser()
	err := p.Parse()

	return &App{}, err
}

func (app *App) Run() error {
	w := worker.NewWorker()
	return w.DoTask()
}
