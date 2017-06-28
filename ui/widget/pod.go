package widget

import (
	"github.com/boz/kcache"
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/backend/nsname"
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/controller"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/view"
	"github.com/gdamore/tcell/views"
)

func NewPodTable(ctx elements.Context, ds pod.BaseDatasource) elements.Widget {
	ctx = ctx.New("pod/table")
	content := table.NewWidget(ctx.Env(), view.PodTableColumns())
	handler := controller.NewPodsPostHandler(ctx, view.NewPodTableWriter(content))
	controller.NewPodsController(ctx, ds, handler)
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

	podDS := podsDS.Filter(backend.NSNamesSelector(nsName))

	svcRootDS, err := ctx.Backend().Services()
	if err != nil {
		podDS.Close()
		ctx.Env().LogErr(err, "service backend")
		return nil, err
	}

	//svcDS := svcRootDS.Filter(backend.ServiceSelector(map[string]string{}))
	svcDS := svcRootDS.Filter(kcache.NullFilter())

	//svcDS := svcRootDS.Filter(kcache.NullFilter())

	pdetails := view.NewPodDetails()

	phandler := controller.NewPodsPostHandler(ctx,
		controller.NewPodHandler(ctx.Env(), pdetails))
	controller.NewPodsController(ctx, podDS, phandler)

	controller.NewPodsController(ctx, podDS,
		controller.NewPodHandler(ctx.Env(), newServiceFilterhandler(svcDS)))

	svcTable := NewServiceTable(ctx, svcDS)

	//svcTable := NewServiceTable(ctx, svcRootDS)
	// return svcTable, nil

	layout := views.NewBoxLayout(views.Vertical)

	layout.AddWidget(pdetails, 0.5)
	layout.AddWidget(svcTable, 1)

	widget := elements.NewWidget(ctx, layout)

	return widget, nil
}

type filterable interface {
	Refilter(kcache.Filter)
}

type refilterHandler struct {
	ds filterable
}

func newServiceFilterhandler(ds filterable) controller.PodHandler {
	return &refilterHandler{ds}
}

func (h *refilterHandler) OnInitialize(obj pod.Pod) {
	h.refilter(obj)
}

func (h *refilterHandler) OnCreate(obj pod.Pod) {
	h.refilter(obj)
}

func (h *refilterHandler) OnUpdate(obj pod.Pod) {
	h.refilter(obj)
}

func (h *refilterHandler) OnDelete(obj pod.Pod) {
	filter := backend.ServiceSelector(map[string]string{})
	h.ds.Refilter(filter)
}

func (h *refilterHandler) refilter(obj pod.Pod) {
	filter := backend.ServiceSelector(obj.Resource().GetLabels())
	h.ds.Refilter(filter)
}
