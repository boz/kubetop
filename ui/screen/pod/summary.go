package pod

import (
	"github.com/boz/kcache/nsname"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	pview "github.com/boz/kubetop/ui/view/pod"
)

func newSummary(ctx elements.Context, id string) (elements.NSNameWidget, error) {
	ctx = ctx.New("pod/summary")

	nsName, err := nsname.Parse(id)
	if err != nil {
		ctx.Env().LogErr(err, "invalid id: %v", id)
		return nil, err
	}

	ds, err := newSummaryDS(ctx, nsName)
	if err != nil {
		ctx.Env().LogErr(err, "invalid id: %v", id)
		return nil, err
	}

	// pod summary
	psummary := pview.NewSummary()
	pview.MonitorUnitary(ctx, ds.pods, psummary)

	// container summary
	ctable := table.NewWidget(ctx.Env(), pview.ContainersTableColumns(), false)
	pview.MonitorUnitary(ctx, ds.pods, pview.NewContainersTable(ctable))

	layout := elements.NewHPanes(true)
	layout.Append(psummary)
	layout.Append(elements.AlignRight(ctable))

	widget := elements.NewNSNameWidget(ctx, layout, nsName)

	return widget, nil
}
