package node

import (
	"github.com/boz/kcache/types/node"
	"github.com/boz/kubetop/ui/elements"
	"k8s.io/api/core/v1"
)

func Poster(poster elements.Poster, delegate node.Handler) node.Handler {
	return node.BuildHandler().
		OnInitialize(func(objs []*v1.Node) {
			poster.PostFunc(func() { delegate.OnInitialize(objs) })
		}).
		OnCreate(func(obj *v1.Node) {
			poster.PostFunc(func() { delegate.OnCreate(obj) })
		}).
		OnUpdate(func(obj *v1.Node) {
			poster.PostFunc(func() { delegate.OnUpdate(obj) })
		}).
		OnDelete(func(obj *v1.Node) {
			poster.PostFunc(func() { delegate.OnDelete(obj) })
		}).Create()
}

func Monitor(
	ctx elements.Context, publisher node.Publisher, handler node.Handler) {
	monitor := node.NewMonitor(publisher, Poster(ctx, handler))
	ctx.AlsoClose(monitor)
}

func MonitorUnitary(
	ctx elements.Context, publisher node.Publisher, handler node.UnitaryHandler) {
	Monitor(ctx, publisher, node.ToUnitary(ctx.Env().Logutil(), handler))
}
