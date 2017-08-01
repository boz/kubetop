package widget

import (
	"github.com/boz/kcache/types/event"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	uiutil "github.com/boz/kubetop/ui/util"
	eview "github.com/boz/kubetop/ui/view/event"
)

func NewEventTable(ctx elements.Context, ds event.Publisher) elements.Widget {
	ctx = ctx.New("event/table")
	content := table.NewWidget(ctx.Env(), eview.TableColumns())

	ctx.AlsoClose(event.NewMonitor(ds,
		uiutil.EventsPoster(ctx, eview.NewTable(content))))

	return elements.NewWidget(ctx, content)
}
