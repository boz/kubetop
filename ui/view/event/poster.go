package table

import (
	"github.com/boz/kcache/types/event"
	"github.com/boz/kubetop/ui/elements"
	"k8s.io/api/core/v1"
)

func Poster(poster elements.Poster, delegate event.Handler) event.Handler {
	return event.BuildHandler().
		OnInitialize(func(objs []*v1.Event) {
			poster.PostFunc(func() { delegate.OnInitialize(objs) })
		}).
		OnCreate(func(obj *v1.Event) {
			poster.PostFunc(func() { delegate.OnCreate(obj) })
		}).
		OnUpdate(func(obj *v1.Event) {
			poster.PostFunc(func() { delegate.OnUpdate(obj) })
		}).
		OnDelete(func(obj *v1.Event) {
			poster.PostFunc(func() { delegate.OnDelete(obj) })
		}).Create()
}

func Monitor(
	ctx elements.Context, publisher event.Publisher, handler event.Handler) {
	monitor, err := event.NewMonitor(publisher, Poster(ctx, handler))
	if err != nil {
		ctx.Env().LogErr(err, "event.NewMonitor")
		return
	}
	ctx.AlsoClose(monitor)
}

func MonitorUnitary(
	ctx elements.Context, publisher event.Publisher, handler event.UnitaryHandler) {
	Monitor(ctx, publisher, event.ToUnitary(ctx.Env().Logutil(), handler))
}
