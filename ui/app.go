package ui

import (
	"github.com/boz/kubetop/backend"
	"github.com/gdamore/tcell/views"
)

type App struct {
	tapp *views.Application

	main *mainWidget

	backend backend.Backend

	stopch chan bool
	donech chan bool
}

func NewApp(backend backend.Backend) *App {
	stopch := make(chan bool, 1)

	tapp := &views.Application{}

	main := newMainWidget(stopch)

	tapp.SetRootWidget(main)

	return &App{
		tapp:   tapp,
		main:   main,
		stopch: stopch,
		donech: make(chan bool),
	}
}

func (app *App) Stop() {
	select {
	case app.stopch <- true:
	case <-app.donech:
	}
}

func (app *App) Run() error {
	defer close(app.donech)
	go func() {
		select {
		case <-app.stopch:
			app.tapp.Quit()
		case <-app.donech:
		}
	}()
	return app.tapp.Run()
}
