package widget

import (
	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/nsname"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	uiutil "github.com/boz/kubetop/ui/util"
	pview "github.com/boz/kubetop/ui/view/pod"
	"k8s.io/api/core/v1"
)

func NewPodTable(ctx elements.Context, ds pod.Publisher) elements.Widget {
	ctx = ctx.New("pod/table")
	content := table.NewWidget(ctx.Env(), pview.TableColumns())
	handler := uiutil.PodsPoster(ctx, pview.NewTable(content))
	ctx.AlsoClose(pod.NewMonitor(ds, handler))
	return elements.NewWidget(ctx, content)
}

func NewPodSummary(ctx elements.Context, id string) (elements.Widget, error) {
	ctx = ctx.New("pod/summary")

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

	// pod summary
	psummary := pview.NewSummary()
	pod.NewMonitor(podController,
		uiutil.PodsPoster(ctx,
			pod.ToUnitary(ctx.Env().Logutil(), psummary)))

	// container summary

	ctable := table.NewWidget(ctx.Env(), pview.ContainersTableColumns())
	pod.NewMonitor(podController,
		uiutil.PodsPoster(ctx,
			pod.ToUnitary(ctx.Env().Logutil(), pview.NewContainersTable(ctable))))

	layout := elements.NewHPanes()
	layout.PushBackWidget(psummary)
	layout.PushBackWidget(ctable)

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
