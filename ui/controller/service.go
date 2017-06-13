package controller

import (
	"github.com/boz/kcache"
	"github.com/boz/kubetop/backend/service"
	"github.com/boz/kubetop/ui/elements"
)

type ServiceHandler interface {
	OnInitialize([]service.Service)
	OnCreate(service.Service)
	OnUpdate(service.Service)
	OnDelete(service.Service)
}

type servicePostHandler struct {
	poster  elements.Poster
	handler ServiceHandler
}

func NewServicePostHandler(poster elements.Poster, handler ServiceHandler) ServiceHandler {
	return &servicePostHandler{poster, handler}
}

func (p *servicePostHandler) OnInitialize(objs []service.Service) {
	p.poster.PostFunc(func() { p.handler.OnInitialize(objs) })
}

func (p *servicePostHandler) OnCreate(obj service.Service) {
	p.poster.PostFunc(func() { p.handler.OnCreate(obj) })
}

func (p *servicePostHandler) OnUpdate(obj service.Service) {
	p.poster.PostFunc(func() { p.handler.OnUpdate(obj) })
}

func (p *servicePostHandler) OnDelete(obj service.Service) {
	p.poster.PostFunc(func() { p.handler.OnDelete(obj) })
}

type ServiceController interface {
}

type serviceController struct {
	sub     service.Subscription
	handler ServiceHandler
	ctx     elements.Context
}

func NewServiceController(ctx elements.Context, ds service.BaseDatasource, handler ServiceHandler) ServiceController {
	controller := &serviceController{ds.Subscribe(kcache.NullFilter()), handler, ctx}
	go controller.run()
	return controller
}

func (c *serviceController) run() {
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
				c.handler.OnCreate(ev.Resource())
			case kcache.EventTypeDelete:
				c.handler.OnDelete(ev.Resource())
			}
		}
	}
}
