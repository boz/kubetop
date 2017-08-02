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
		ctx.Env().LogErr(err, "pod ds")
		return nil, err
	}

	pods, err := pbase.CloneWithFilter(filter.NSName(id))
	if err != nil {
		ctx.Env().LogErr(err, "pod ds")
		return nil, err
	}
	ctx.AlsoClose(pods)

	return &podSummaryDS{pods}, nil
}

type podIndexDS struct {
	pods pod.Controller
}

func newIndexDS(ctx elements.Context) (*podIndexDS, error) {
	pbase, err := ctx.Backend().Pods()
	if err != nil {
		return nil, err
	}
	pods, err := pbase.Clone()
	if err != nil {
		return nil, err
	}
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
	pods, err := podsBase.CloneWithFilter(filter.NSName(id))
	if err != nil {
		return nil, err
	}
	ctx.AlsoClose(pods)

	servicesBase, err := ctx.Backend().Services()
	if err != nil {
		return nil, err
	}
	services, err := servicesBase.CloneWithFilter(filter.All())
	if err != nil {
		return nil, err
	}
	ctx.AlsoClose(services)

	eventsBase, err := ctx.Backend().Events()
	if err != nil {
		return nil, err
	}
	events, err := eventsBase.CloneWithFilter(filter.All())
	if err != nil {
		return nil, err
	}
	ctx.AlsoClose(events)

	m, err := pod.NewMonitor(pods,
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

	if err != nil {
		return nil, err
	}
	ctx.AlsoClose(m)

	return &showDS{pods, services, events}, nil
}
