package controller

import (
	"github.com/boz/kcache"
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/elements"
)

type PodsHandler interface {
	OnInitialize([]pod.Pod)
	OnCreate(pod.Pod)
	OnUpdate(pod.Pod)
	OnDelete(pod.Pod)
}

type podsPostHandler struct {
	poster  elements.Poster
	handler PodsHandler
}

func NewPodsPostHandler(poster elements.Poster, handler PodsHandler) PodsHandler {
	return &podsPostHandler{poster, handler}
}

func (p *podsPostHandler) OnInitialize(objs []pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnInitialize(objs) })
}

func (p *podsPostHandler) OnCreate(obj pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnCreate(obj) })
}

func (p *podsPostHandler) OnUpdate(obj pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnUpdate(obj) })
}

func (p *podsPostHandler) OnDelete(obj pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnDelete(obj) })
}

type PodController interface {
}

type podsController struct {
	sub     pod.Subscription
	handler PodsHandler
	ctx     elements.Context
}

func NewPodsController(ctx elements.Context, ds pod.BaseDatasource, handler PodsHandler) PodController {
	controller := &podsController{ds.Subscribe(kcache.NullFilter()), handler, ctx}
	go controller.run()
	return controller
}

func (c *podsController) run() {
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
				c.handler.OnUpdate(ev.Resource())
			case kcache.EventTypeDelete:
				c.handler.OnDelete(ev.Resource())
			}
		}
	}
}
