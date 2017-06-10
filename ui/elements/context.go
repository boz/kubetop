package elements

import (
	lifecycle "github.com/boz/go-lifecycle"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/util"
)

type Context interface {
	Env() util.Env
	Backend() backend.Backend

	Close()
	Closed() <-chan struct{}
	WatchChannel(<-chan struct{})

	New(string) Context
	NewWithID(string) Context
	SetEnv(env util.Env)

	PostFunc(fn func())
}

type Poster interface {
	PostFunc(fn func())
}

type context struct {
	poster  Poster
	backend backend.Backend
	env     util.Env
	lc      lifecycle.Lifecycle
}

func NewContext(env util.Env, backend backend.Backend, poster Poster) Context {

	lc := lifecycle.New()
	go func() {
		defer lc.ShutdownCompleted()
		defer lc.ShutdownInitiated()
		<-lc.ShutdownRequest()
	}()

	return &context{
		poster:  poster,
		backend: backend,
		env:     env,
		lc:      lc,
	}
}

func (p *context) New(name string) Context {
	new := NewContext(p.env.ForComponent(name), p.backend, p.poster)
	go new.WatchChannel(p.Closed())
	return new
}

func (p *context) NewWithID(name string) Context {
	new := p.New(name)
	new.SetEnv(new.Env().WithID())
	return new
}

func (p *context) SetEnv(env util.Env) {
	p.env = env
}

func (p *context) Close() {
	p.lc.Shutdown()
}

func (p *context) Closed() <-chan struct{} {
	return p.lc.Done()
}

func (p *context) WatchChannel(ch <-chan struct{}) {
	go p.lc.WatchChannel(ch)
}

func (p *context) PostFunc(fn func()) {
	p.poster.PostFunc(fn)
}

func (p *context) Env() util.Env {
	return p.env
}

func (p *context) Backend() backend.Backend {
	return p.backend
}
