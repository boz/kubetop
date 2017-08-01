package pod

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/nsname"
	"github.com/boz/kcache/types/event"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	uiutil "github.com/boz/kubetop/ui/util"
	eview "github.com/boz/kubetop/ui/view/event"
	pview "github.com/boz/kubetop/ui/view/pod"
	sview "github.com/boz/kubetop/ui/view/service"
	"github.com/gdamore/tcell/views"
)

type showScreen struct {
	layout elements.Panes
	ctx    elements.Context
	views.WidgetWatchers
}

func NewShow(ctx elements.Context, req elements.NSNameRequest) (elements.Screen, error) {
	ctx = ctx.New("pod/show")

	ds, err := newShowDS(ctx, req.NSName())
	if err != nil {
		ctx.Env().LogErr(err, "opening data source")
		ctx.Close()
		return nil, err
	}

	layout := elements.NewVPanes(true)

	containerw := showContainersPanel(ctx, ds)
	layout.PushBackWidget(containerw)

	servicesw := showServicesPanel(ctx, ds)
	layout.PushBackWidget(servicesw)

	eventsw := showEventsPanel(ctx, ds)
	layout.PushBackWidget(eventsw)

	return elements.NewScreen(ctx, req, fmt.Sprintf("Pod %v", req.NSName()), layout), nil
}

func showContainersPanel(ctx elements.Context, ds showDS) views.Widget {
	ctx = ctx.New("pod/show/containers")
	ctable := table.NewWidget(ctx.Env(), pview.ContainersTableColumns(), true)
	pod.NewMonitor(ds.pods,
		uiutil.PodsPoster(ctx,
			pod.ToUnitary(ctx.Env().Logutil(), pview.NewContainersTable(ctable))))
	return newPanel("Containers", ctable)
}

func showServicesPanel(ctx elements.Context, ds showDS) views.Widget {
	ctx = ctx.New("pod/show/services")
	content := table.NewWidget(ctx.Env(), sview.TableColumns(), true)
	service.NewMonitor(ds.services,
		uiutil.ServicesPoster(ctx,
			sview.NewTable(content)))
	return newPanel("Services", content)
}

func showEventsPanel(ctx elements.Context, ds showDS) views.Widget {
	ctx = ctx.New("pod/show/events")
	content := table.NewWidget(ctx.Env(), eview.TableColumns(), true)
	event.NewMonitor(ds.events,
		uiutil.EventsPoster(ctx,
			eview.NewTable(content)))
	return newPanel("Events", content)
}

func newPanel(title string, content views.Widget) views.Widget {
	titlew := views.NewTextBar()
	titlew.SetCenter(title, theme.Base)
	layout := elements.NewVPanes(false)
	layout.PushBackWidget(titlew)
	layout.PushBackWidget(content)
	return layout
}

type showDS struct {
	pods     pod.Publisher
	services service.Publisher
	events   event.Publisher
}

func newShowDS(ctx elements.Context, id nsname.NSName) (showDS, error) {
	podsBase, err := ctx.Backend().Pods()
	if err != nil {
		return showDS{}, err
	}
	pods := podsBase.CloneWithFilter(filter.NSName(id))
	ctx.AlsoClose(pods)

	servicesBase, err := ctx.Backend().Services()
	if err != nil {
		return showDS{}, err
	}
	services := servicesBase.CloneWithFilter(filter.All())
	ctx.AlsoClose(services)

	eventsBase, err := ctx.Backend().Events()
	if err != nil {
		return showDS{}, err
	}
	events := eventsBase.CloneWithFilter(filter.All())
	ctx.AlsoClose(events)

	pod.NewMonitor(pods,
		pod.ToUnitary(ctx.Env().Logutil(),
			pod.BuildUnitaryHandler().
				OnInitialize(func(obj *v1.Pod) {
					events.Refilter(
						event.InvolvedFilter("Pod", obj.GetNamespace(), obj.GetName()))
					services.Refilter(service.SelectorMatchFilter(obj.GetLabels()))
				}).
				OnCreate(func(obj *v1.Pod) {
					events.Refilter(
						event.InvolvedFilter("Pod", obj.GetNamespace(), obj.GetName()))
					services.Refilter(service.SelectorMatchFilter(obj.GetLabels()))
				}).
				OnUpdate(func(obj *v1.Pod) {
					events.Refilter(
						event.InvolvedFilter("Pod", obj.GetNamespace(), obj.GetName()))
					services.Refilter(service.SelectorMatchFilter(obj.GetLabels()))
				}).
				OnDelete(func(obj *v1.Pod) {
					events.Refilter(filter.All())
					services.Refilter(filter.All())
				}).Create()))

	return showDS{pods, services, events}, nil
}
