package widget

import (
	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/nsname"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	uiutil "github.com/boz/kubetop/ui/util"
	"github.com/boz/kubetop/ui/view"
	"github.com/gdamore/tcell/views"
	"k8s.io/api/core/v1"
)

func NewPodTable(ctx elements.Context, ds pod.Publisher) elements.Widget {
	ctx = ctx.New("pod/table")
	content := table.NewWidget(ctx.Env(), view.PodTableColumns())
	handler := uiutil.PodsPoster(ctx, view.NewPodTableWriter(content))
	ctx.AlsoClose(pod.NewMonitor(ds, handler))
	return elements.NewWidget(ctx, content)
}

func NewPodDetails(ctx elements.Context, id string) (elements.Widget, error) {
	ctx = ctx.New("pod/details")

	nsName, err := nsname.Parse(id)
	if err != nil {
		ctx.Env().LogErr(err, "invalid id: %v", id)
		return nil, err
	}

	podsBase, err := ctx.Backend().Pods()
	if err != nil {
		ctx.Env().LogErr(err, "pod backend")
		return nil, err
	}

	podController := podsBase.CloneWithFilter(filter.NSName(nsName))
	ctx.AlsoClose(podController)

	// pod details
	pdetails := view.NewPodDetails()
	uiutil.PodsPoster(ctx,
		pod.ToUnitary(ctx.Env().Logutil(), pdetails))

	svcBase, err := ctx.Backend().Services()
	if err != nil {
		ctx.Env().LogErr(err, "service backend")
		ctx.Close()
		return nil, err
	}

	svcController := svcBase.CloneWithFilter(filter.All())
	ctx.AlsoClose(svcController)

	// keep svcController up to date with pod
	pod.NewMonitor(podController,
		pod.ToUnitary(ctx.Env().Logutil(),
			podServicesHandler(svcController)))

	// display services matching pod
	svcTable := NewServiceTable(ctx, svcController)

	layout := views.NewBoxLayout(views.Vertical)

	layout.AddWidget(pdetails, 0.5)
	layout.AddWidget(svcTable, 1)

	widget := elements.NewWidget(ctx, layout)

	return widget, nil
}

func podServicesHandler(target filterable) pod.UnitaryHandler {
	return pod.BuildUnitaryHandler().
		OnInitialize(func(obj *v1.Pod) {
			target.Refilter(service.SelectorMatchFilter(obj.GetLabels()))
		}).
		OnCreate(func(obj *v1.Pod) {
			target.Refilter(service.SelectorMatchFilter(obj.GetLabels()))
		}).
		OnUpdate(func(obj *v1.Pod) {
			target.Refilter(service.SelectorMatchFilter(obj.GetLabels()))
		}).
		OnDelete(func(obj *v1.Pod) {
			target.Refilter(filter.All())
		}).Create()
}

type filterable interface {
	Refilter(filter.Filter)
}
