package widget

import (
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	uiutil "github.com/boz/kubetop/ui/util"
	"github.com/boz/kubetop/ui/view"
)

func NewServiceTable(ctx elements.Context, ds service.Publisher) elements.Widget {
	ctx = ctx.New("service/table")
	content := table.NewWidget(ctx.Env(), view.ServiceTableColumns())

	ctx.AlsoClose(service.NewMonitor(ds,
		uiutil.ServicesPoster(ctx, view.NewServiceTableWriter(content))))

	return elements.NewWidget(ctx, content)
}
