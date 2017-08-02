package service

import (
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"k8s.io/api/core/v1"
)

func Poster(poster elements.Poster, delegate service.Handler) service.Handler {
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

func Monitor(
	ctx elements.Context, publisher service.Publisher, handler service.Handler) {
	monitor := service.NewMonitor(publisher, Poster(ctx, handler))
	ctx.AlsoClose(monitor)
}

func MonitorUnitary(
	ctx elements.Context, publisher service.Publisher, handler service.UnitaryHandler) {
	Monitor(ctx, publisher, service.ToUnitary(ctx.Env().Logutil(), handler))
}
