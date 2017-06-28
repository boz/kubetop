package controller

import (
	"github.com/boz/kcache"
	"github.com/boz/kubetop/backend/service"
	"github.com/boz/kubetop/ui/elements"
)

type ServicesHandler interface {
	OnInitialize([]service.Service)
	OnCreate(service.Service)
	OnUpdate(service.Service)
	OnDelete(service.Service)
}

type servicesPostHandler struct {
	poster  elements.Poster
	handler ServicesHandler
}

func NewServicesPostHandler(poster elements.Poster, handler ServicesHandler) ServicesHandler {
	return &servicesPostHandler{poster, handler}
}

func (p *servicesPostHandler) OnInitialize(objs []service.Service) {
	p.poster.PostFunc(func() { p.handler.OnInitialize(objs) })
}

func (p *servicesPostHandler) OnCreate(obj service.Service) {
	p.poster.PostFunc(func() { p.handler.OnCreate(obj) })
}

func (p *servicesPostHandler) OnUpdate(obj service.Service) {
	p.poster.PostFunc(func() { p.handler.OnUpdate(obj) })
}

func (p *servicesPostHandler) OnDelete(obj service.Service) {
	p.poster.PostFunc(func() { p.handler.OnDelete(obj) })
}

type ServicesController interface {
}

type servicesController struct {
	sub     service.Subscription
	handler ServicesHandler
	ctx     elements.Context
}

func NewServiceController(ctx elements.Context, ds service.BaseDatasource, handler ServicesHandler) ServicesController {
	controller := &servicesController{ds.Subscribe(kcache.NullFilter()), handler, ctx}
	go controller.run()
	return controller
}

func (c *servicesController) run() {
	defer c.sub.Close()

	readych := c.sub.Ready()
	for {
		select {
		case <-c.ctx.Closed():
			return
		case <-c.sub.Closed():
			return
		case <-readych:
			objs, _ := c.sub.List()
			c.ctx.Env().Log().Debugf("%v services", len(objs))
			c.handler.OnInitialize(objs)
			readych = nil
		case ev, ok := <-c.sub.Events():
			if !ok {
				return
			}
			switch ev.Type() {
			case kcache.EventTypeCreate:
				c.handler.OnCreate(ev.Resource())
			case kcache.EventTypeUpdate:
				c.handler.OnUpdate(ev.Resource())
			case kcache.EventTypeDelete:
				c.handler.OnDelete(ev.Resource())
			}
		}
	}
}
