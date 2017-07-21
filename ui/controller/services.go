package controller

import (
	"github.com/boz/kcache"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"k8s.io/api/core/v1"
)

type ServicesHandler interface {
	OnInitialize([]*v1.Service)
	OnCreate(*v1.Service)
	OnUpdate(*v1.Service)
	OnDelete(*v1.Service)
}

type servicesPostHandler struct {
	poster  elements.Poster
	handler ServicesHandler
}

func NewServicesPostHandler(poster elements.Poster, handler ServicesHandler) ServicesHandler {
	return &servicesPostHandler{poster, handler}
}

func (p *servicesPostHandler) OnInitialize(objs []*v1.Service) {
	p.poster.PostFunc(func() { p.handler.OnInitialize(objs) })
}

func (p *servicesPostHandler) OnCreate(obj *v1.Service) {
	p.poster.PostFunc(func() { p.handler.OnCreate(obj) })
}

func (p *servicesPostHandler) OnUpdate(obj *v1.Service) {
	p.poster.PostFunc(func() { p.handler.OnUpdate(obj) })
}

func (p *servicesPostHandler) OnDelete(obj *v1.Service) {
	p.poster.PostFunc(func() { p.handler.OnDelete(obj) })
}

type ServicesController interface {
}

type servicesController struct {
	sub     service.Subscription
	handler ServicesHandler
	ctx     elements.Context
}

func NewServiceController(ctx elements.Context, ds service.Publisher, handler ServicesHandler) ServicesController {
	controller := &servicesController{ds.Subscribe(), handler, ctx}
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
		case <-c.sub.Done():
			return
		case <-readych:
			objs, _ := c.sub.Cache().List()
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
