package controller

import (
	"github.com/boz/kcache"
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/elements"
)

type PodHandler interface {
	OnInitialize([]pod.Pod)
	OnCreate(pod.Pod)
	OnUpdate(pod.Pod)
	OnDelete(pod.Pod)
}

type podPostHandler struct {
	poster  elements.Poster
	handler PodHandler
}

func NewPodPostHandler(poster elements.Poster, handler PodHandler) PodHandler {
	return &podPostHandler{poster, handler}
}

func (p *podPostHandler) OnInitialize(objs []pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnInitialize(objs) })
}

func (p *podPostHandler) OnCreate(obj pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnCreate(obj) })
}

func (p *podPostHandler) OnUpdate(obj pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnUpdate(obj) })
}

func (p *podPostHandler) OnDelete(obj pod.Pod) {
	p.poster.PostFunc(func() { p.handler.OnDelete(obj) })
}

type PodController interface {
}

type podController struct {
	sub     pod.Subscription
	handler PodHandler
	ctx     elements.Context
}

func NewPodController(ctx elements.Context, ds pod.BaseDatasource, handler PodHandler) PodController {
	controller := &podController{ds.Subscribe(kcache.NullFilter()), handler, ctx}
	go controller.run()
	return controller
}

func (c *podController) run() {
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
