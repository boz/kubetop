package elements

import (
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/util"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Presenter interface {
	views.Widget

	Env() util.Env
	Backend() backend.Backend

	Close()
	Closed() <-chan struct{}

	New(string, views.Widget) Presenter
	NewWithID(string, views.Widget) Presenter
	SetEnv(env util.Env)

	PostFunc(fn func())
}

type isPresentable interface {
	SetPresenter(Presenter)
	Presenter() Presenter
}

type funcPoster interface {
	PostFunc(fn func())
}

type Presentable struct {
	presenter Presenter
}

func ClosePresenter(w views.Widget) {
	if w, ok := w.(isPresentable); ok {
		w.Presenter().Close()
	}
}

func (p *Presentable) SetPresenter(presenter Presenter) {
	p.presenter = presenter
}

func (p *Presentable) Presenter() Presenter {
	return p.presenter
}

func (p *Presentable) Env() util.Env {
	return p.presenter.Env()
}

func (p *Presentable) Backend() backend.Backend {
	return p.presenter.Backend()
}

type presenter struct {
	content views.Widget
	poster  funcPoster

	backend backend.Backend
	env     util.Env

	closedch chan struct{}
}

func NewPresenter(env util.Env, backend backend.Backend, poster funcPoster, content views.Widget) Presenter {
	new := &presenter{
		content:  content,
		poster:   poster,
		backend:  backend,
		env:      env,
		closedch: make(chan struct{}),
	}
	if content, ok := content.(isPresentable); ok {
		content.SetPresenter(new)
	}
	return new
}

func (p *presenter) New(name string, content views.Widget) Presenter {
	new := &presenter{
		content:  content,
		poster:   p.poster,
		backend:  p.backend,
		env:      p.env.ForComponent(name),
		closedch: make(chan struct{}),
	}
	if content, ok := content.(isPresentable); ok {
		content.SetPresenter(new)
	}
	return new
}

func (p *presenter) NewWithID(name string, content views.Widget) Presenter {
	new := p.New(name, content)
	new.SetEnv(new.Env().WithID())
	return new
}

func (p *presenter) SetEnv(env util.Env) {
	p.env = env
}

func (p *presenter) Close() {
	close(p.closedch)
}

func (p *presenter) Closed() <-chan struct{} {
	return p.closedch
}

func (p *presenter) PostFunc(fn func()) {
	p.poster.PostFunc(fn)
}

func (p *presenter) Env() util.Env {
	return p.env
}

func (p *presenter) Backend() backend.Backend {
	return p.backend
}

// views.Widget methods

func (p *presenter) Draw()                              { p.content.Draw() }
func (p *presenter) Resize()                            { p.content.Resize() }
func (p *presenter) HandleEvent(ev tcell.Event) bool    { return p.content.HandleEvent(ev) }
func (p *presenter) SetView(view views.View)            { p.content.SetView(view) }
func (p *presenter) Size() (int, int)                   { return p.content.Size() }
func (p *presenter) Watch(handler tcell.EventHandler)   { p.content.Watch(handler) }
func (p *presenter) Unwatch(handler tcell.EventHandler) { p.content.Unwatch(handler) }
