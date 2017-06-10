package ui

import (
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/util"
	"github.com/gdamore/tcell/views"
)

type App struct {
	tapp *views.Application

	main views.Widget

	stopch chan bool
	donech chan bool

	ctx elements.Context
}

func NewApp(env util.Env, backend backend.Backend) *App {
	env = env.ForComponent("ui/app")

	stopch := make(chan bool, 1)
	tapp := &views.Application{}

	ctx := elements.NewContext(env, backend, tapp)
	main := newMainWidget(ctx, stopch)

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
