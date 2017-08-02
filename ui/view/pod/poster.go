package pod

import (
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements"
	"k8s.io/api/core/v1"
)

func Poster(poster elements.Poster, delegate pod.Handler) pod.Handler {
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

func Monitor(
	ctx elements.Context, publisher pod.Publisher, handler pod.Handler) {
	monitor, err := pod.NewMonitor(publisher, Poster(ctx, handler))
	if err != nil {
		ctx.Env().LogErr(err, "pod.NewMonitor")
		return
	}
	ctx.AlsoClose(monitor)
}

func MonitorUnitary(
	ctx elements.Context, publisher pod.Publisher, handler pod.UnitaryHandler) {
	Monitor(ctx, publisher, pod.ToUnitary(ctx.Env().Logutil(), handler))
}
