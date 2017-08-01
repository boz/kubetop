package widget

import (
	"github.com/boz/kcache/types/node"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	uiutil "github.com/boz/kubetop/ui/util"
	nview "github.com/boz/kubetop/ui/view/node"
)

func NewNodeTable(ctx elements.Context, ds node.Publisher) elements.Widget {
	ctx = ctx.New("node/table")
	content := table.NewWidget(ctx.Env(), nview.TableColumns(), true)

	ctx.AlsoClose(node.NewMonitor(ds,
		uiutil.NodesPoster(ctx, nview.NewTable(content))))

	return elements.NewWidget(ctx, content)
}
