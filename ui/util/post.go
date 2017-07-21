package util

import (
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"k8s.io/api/core/v1"
)

func PodsPoster(poster elements.Poster, delegate pod.Handler) pod.Handler {
	return pod.BuildHandler().
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

func ServicesPoster(poster elements.Poster, delegate service.Handler) service.Handler {
	return service.BuildHandler().
		OnInitialize(func(objs []*v1.Service) {
			poster.PostFunc(func() { delegate.OnInitialize(objs) })
		}).
		OnCreate(func(obj *v1.Service) {
			poster.PostFunc(func() { delegate.OnCreate(obj) })
		}).
		OnUpdate(func(obj *v1.Service) {
			poster.PostFunc(func() { delegate.OnUpdate(obj) })
		}).
		OnDelete(func(obj *v1.Service) {
			poster.PostFunc(func() { delegate.OnDelete(obj) })
		}).Create()
}
