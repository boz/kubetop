package event

import (
	"github.com/boz/kcache/types/event"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	eview "github.com/boz/kubetop/ui/view/event"
)

func newIndexTable(ctx elements.Context, ds event.Publisher) elements.Widget {
	ctx = ctx.New("event/table")
	content := table.NewWidget(ctx.Env(), eview.TableColumns(), true)
	eview.Monitor(ctx, ds, eview.NewTable(content))
	return elements.NewWidget(ctx, content)
}
