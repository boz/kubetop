package widget

import (
	"github.com/boz/kcache/types/node"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	uiutil "github.com/boz/kubetop/ui/util"
	"github.com/boz/kubetop/ui/view"
)

func NewNodeTable(ctx elements.Context, ds node.Publisher) elements.Widget {
	ctx = ctx.New("node/table")
	content := table.NewWidget(ctx.Env(), view.NodeTableColumns())

	ctx.AlsoClose(node.NewMonitor(ds,
		uiutil.NodesPoster(ctx, view.NewNodeTableWriter(content))))

	return elements.NewWidget(ctx, content)
}
