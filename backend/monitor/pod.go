package monitor

import (
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/util"
	"k8s.io/api/core/v1"
)

type PodHandler interface {
	OnInitialize(*v1.Pod)
	OnCreate(*v1.Pod)
	OnUpdate(*v1.Pod)
	OnDelete(*v1.Pod)
}

func NewPodHandler(env util.Env, delegate PodHandler) pod.Handler {
	return pod.NewHandlerBuilder().
		OnInitialize(func(objs []*v1.Pod) {
			if count := len(objs); count > 1 {
				env.Log().Warnf("initialized with invalid count: %v", count)
				return
			}
			if count := len(objs); count == 0 {
				env.Log().Debugf("initialized with empty result, ignoring")
				return
			}
			delegate.OnInitialize(objs[0])
		}).
		OnCreate(func(obj *v1.Pod) {
			delegate.OnCreate(obj)
		}).
		OnUpdate(func(obj *v1.Pod) {
			delegate.OnUpdate(obj)
		}).
		OnDelete(func(obj *v1.Pod) {
			delegate.OnDelete(obj)
		}).Create()
}

func NewPodsPostHandler(poster elements.Poster, delegate pod.Handler) pod.Handler {
	return pod.NewHandlerBuilder().
		OnInitialize(func(objs []*v1.Pod) {
			poster.PostFunc(func() { delegate.OnInitialize(objs) })
		}).
		OnCreate(func(obj *v1.Pod) {
			poster.PostFunc(func() { delegate.OnCreate(obj) })
		}).
		OnUpdate(func(obj *v1.Pod) {
			poster.PostFunc(func() { delegate.OnUpdate(obj) })
		}).
		OnDelete(func(obj *v1.Pod) {
			poster.PostFunc(func() { delegate.OnDelete(obj) })
		}).Create()
}
