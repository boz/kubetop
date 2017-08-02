package pod

import (
	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/nsname"
	"github.com/boz/kcache/types/event"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"k8s.io/api/core/v1"
)

type podSummaryDS struct {
	pods pod.FilterController
}

func newSummaryDS(ctx elements.Context, id nsname.NSName) (*podSummaryDS, error) {
	pbase, err := ctx.Backend().Pods()

	if err != nil {
		ctx.Env().LogErr(err, "pod datasource")
		return nil, err
	}

	pods := pbase.CloneWithFilter(filter.NSName(id))
	ctx.AlsoClose(pods)

	return &podSummaryDS{pods}, nil
}

type podIndexDS struct {
	pods pod.Controller
}

func newIndexDS(ctx elements.Context) (*podIndexDS, error) {
	pbase, err := ctx.Backend().Pods()
	if err != nil {
		ctx.Env().LogErr(err, "pod datasource")
		return nil, err
	}
	pods := pbase.Clone()
	ctx.AlsoClose(pods)
	return &podIndexDS{pods}, nil
}

type showDS struct {
	pods     pod.Publisher
	services service.Publisher
	events   event.Publisher
}

func newShowDS(ctx elements.Context, id nsname.NSName) (*showDS, error) {
	podsBase, err := ctx.Backend().Pods()
	if err != nil {
		return nil, err
	}
	pods := podsBase.CloneWithFilter(filter.NSName(id))
	ctx.AlsoClose(pods)

	servicesBase, err := ctx.Backend().Services()
	if err != nil {
		return nil, err
	}
	services := servicesBase.CloneWithFilter(filter.All())
	ctx.AlsoClose(services)

	eventsBase, err := ctx.Backend().Events()
	if err != nil {
		return nil, err
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

	return &showDS{pods, services, events}, nil
}
