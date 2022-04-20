package app

import (
	"loki/internal/loki/config"
	"loki/internal/loki/parser"
	"loki/internal/loki/worker"
)

type App struct {
	config config.Config
}

func NewApp() (*App, error) {
	p := parser.NewParser()
	err := p.Parse()

	return &App{config: *config.AppConfig}, err
}

func (app *App) Run() error {
	w := worker.NewWorker(app.config)
	if err := w.DoTask(); err != nil {
		return err
	}

	return nil
}
