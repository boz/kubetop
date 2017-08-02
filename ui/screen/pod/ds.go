package pod

import (
	"github.com/boz/kcache/filter"
	"github.com/boz/kcache/nsname"
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements"
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
