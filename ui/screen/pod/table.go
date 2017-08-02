package pod

import (
	"github.com/boz/kcache/types/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	pview "github.com/boz/kubetop/ui/view/pod"
)

func newIndexTable(ctx elements.Context, ds pod.Publisher) elements.Widget {
	ctx = ctx.New("pod/index#table")
	content := table.NewWidget(ctx.Env(), pview.TableColumns(), true)
	pview.Monitor(ctx, ds, pview.NewTable(content))
	return elements.NewWidget(ctx, content)
}
