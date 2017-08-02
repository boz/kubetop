package node

import (
	"github.com/boz/kcache/types/node"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	nview "github.com/boz/kubetop/ui/view/node"
)

func newIndexTable(ctx elements.Context, ds node.Publisher) elements.Widget {
	ctx = ctx.New("node/table")
	content := table.NewWidget(ctx.Env(), nview.TableColumns(), true)
	nview.Monitor(ctx, ds, nview.NewTable(content))
	return elements.NewWidget(ctx, content)
}
