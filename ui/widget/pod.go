package widget

import (
	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/nsname"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/backend/monitor"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/view"
	"github.com/gdamore/tcell/views"
	"k8s.io/api/core/v1"
)

func NewPodTable(ctx elements.Context, ds pod.Publisher) elements.Widget {
	ctx = ctx.New("pod/table")
	content := table.NewWidget(ctx.Env(), view.PodTableColumns())
	handler := monitor.NewPodsPostHandler(ctx, view.NewPodTableWriter(content))
	ctx.OnClose(pod.NewMonitor(ds, handler).Close)
	return elements.NewWidget(ctx, content)
}

func NewPodDetails(ctx elements.Context, id string) (elements.Widget, error) {
	ctx = ctx.New("pod/details")

	nsName, err := nsname.Parse(id)
	if err != nil {
		ctx.Env().LogErr(err, "invalid id: %v", id)
		return nil, err
	}

	podsDS, err := ctx.Backend().Pods()
	if err != nil {
		ctx.Env().LogErr(err, "pod backend")
		return nil, err
	}

	podDS := podsDS.CloneWithFilter(filter.NSName(nsName))
	ctx.OnClose(podDS.Close)

	svcRootDS, err := ctx.Backend().Services()
	if err != nil {
		ctx.Close()
		ctx.Env().LogErr(err, "service backend")
		return nil, err
	}

	svcDS := svcRootDS.CloneWithFilter(filter.All())
	ctx.OnClose(svcDS.Close)

	pdetails := view.NewPodDetails()

	phandler := monitor.NewPodsPostHandler(ctx, monitor.NewPodHandler(ctx.Env(), pdetails))
	ctx.OnClose(pod.NewMonitor(podDS, phandler).Close)

	ctx.OnClose(pod.NewMonitor(podDS, monitor.NewPodHandler(ctx.Env(), newServiceFilterhandler(svcDS))).Close)

	svcTable := NewServiceTable(ctx, svcDS)

	layout := views.NewBoxLayout(views.Vertical)

	layout.AddWidget(pdetails, 0.5)
	layout.AddWidget(svcTable, 1)

	widget := elements.NewWidget(ctx, layout)

	return widget, nil
}

type filterable interface {
	Refilter(filter.Filter)
}

type refilterHandler struct {
	ds filterable
}

func newServiceFilterhandler(ds filterable) monitor.PodHandler {
	return &refilterHandler{ds}
}

func (h *refilterHandler) OnInitialize(obj *v1.Pod) {
	h.refilter(obj)
}

func (h *refilterHandler) OnCreate(obj *v1.Pod) {
	h.refilter(obj)
}

func (h *refilterHandler) OnUpdate(obj *v1.Pod) {
	h.refilter(obj)
}

func (h *refilterHandler) OnDelete(obj *v1.Pod) {
	filter := service.SelectorMatchFilter(map[string]string{})
	h.ds.Refilter(filter)
}

func (h *refilterHandler) refilter(obj *v1.Pod) {
	filter := service.SelectorMatchFilter(obj.GetLabels())
	h.ds.Refilter(filter)
}
