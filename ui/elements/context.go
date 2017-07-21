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

	OnClose(fn func())
	AlsoClose(Closeable)

	WatchNavigation(NavWatcher)
	NavigateTo(Request)
}

type NavWatcher interface {
	HandleNavigationRequest(Request)
}

type Poster interface {
	PostFunc(fn func())
}

type Closeable interface {
	Close()
}

type context struct {
	poster      Poster
	navWatchers map[NavWatcher]bool
	backend     backend.Backend
	env         util.Env
	lc          lifecycle.Lifecycle
}

func NewContext(env util.Env, backend backend.Backend, poster Poster) Context {

	lc := lifecycle.New()
	go func() {
		defer lc.ShutdownCompleted()
		defer lc.ShutdownInitiated()
		<-lc.ShutdownRequest()
	}()

	return &context{
		poster:      poster,
		navWatchers: make(map[NavWatcher]bool),
		backend:     backend,
		env:         env,
		lc:          lc,
	}
}

func (p *context) clone(name string) *context {
	lc := lifecycle.New()
	go func() {
		defer lc.ShutdownCompleted()
		defer lc.ShutdownInitiated()
		<-lc.ShutdownRequest()
	}()

	go lc.WatchChannel(p.Closed())

	return &context{
		poster:      p.poster,
		navWatchers: make(map[NavWatcher]bool),
		backend:     p.backend,
		env:         p.env.ForComponent(name),
		lc:          lc,
	}
}

func (p *context) New(name string) Context {
	return p.clone(name)
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

func (p *context) OnClose(fn func()) {
	go func() {
		<-p.lc.Done()
		fn()
	}()
}

func (p *context) AlsoClose(child Closeable) {
	p.OnClose(child.Close)
}

func (p *context) WatchNavigation(nh NavWatcher) {
	p.navWatchers[nh] = true
}

func (p *context) NavigateTo(req Request) {
	for nw, _ := range p.navWatchers {
		nw.HandleNavigationRequest(req)
	}
}
