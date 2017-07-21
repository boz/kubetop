package controller

import (
	"github.com/boz/kcache"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/util"
	"k8s.io/api/core/v1"
)

type PodsHandler interface {
	OnInitialize([]*v1.Pod)
	OnCreate(*v1.Pod)
	OnUpdate(*v1.Pod)
	OnDelete(*v1.Pod)
}

type PodHandler interface {
	OnInitialize(*v1.Pod)
	OnCreate(*v1.Pod)
	OnUpdate(*v1.Pod)
	OnDelete(*v1.Pod)
}

type podsPostHandler struct {
	poster  elements.Poster
	handler PodsHandler
}

func NewPodsPostHandler(poster elements.Poster, handler PodsHandler) PodsHandler {
	return &podsPostHandler{poster, handler}
}

func (p *podsPostHandler) OnInitialize(objs []*v1.Pod) {
	p.poster.PostFunc(func() { p.handler.OnInitialize(objs) })
}

func (p *podsPostHandler) OnCreate(obj *v1.Pod) {
	p.poster.PostFunc(func() { p.handler.OnCreate(obj) })
}

func (p *podsPostHandler) OnUpdate(obj *v1.Pod) {
	p.poster.PostFunc(func() { p.handler.OnUpdate(obj) })
}

func (p *podsPostHandler) OnDelete(obj *v1.Pod) {
	p.poster.PostFunc(func() { p.handler.OnDelete(obj) })
}

type podHandler struct {
	delegate PodHandler
	env      util.Env
}

type PodController interface {
}

type podsController struct {
	sub     pod.Subscription
	handler PodsHandler
	ctx     elements.Context
}

func NewPodsController(ctx elements.Context, ds pod.Publisher, handler PodsHandler) PodController {
	controller := &podsController{ds.Subscribe(), handler, ctx}
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

func NewPodHandler(env util.Env, delegate PodHandler) PodsHandler {
	return &podHandler{delegate, env}
}

func (p *podHandler) OnInitialize(objs []*v1.Pod) {

	if count := len(objs); count > 1 {
		p.env.Log().Warnf("initialized with invalid count: %v", count)
		return
	}

	if count := len(objs); count == 0 {
		p.env.Log().Debugf("initialized with empty result, ignoring")
		return
	}

	p.delegate.OnInitialize(objs[0])
}

func (p *podHandler) OnCreate(obj *v1.Pod) {
	p.delegate.OnCreate(obj)
}

func (p *podHandler) OnUpdate(obj *v1.Pod) {
	p.delegate.OnUpdate(obj)
}

func (p *podHandler) OnDelete(obj *v1.Pod) {
	p.delegate.OnDelete(obj)
}
