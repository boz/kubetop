package widget

import (
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	uiutil "github.com/boz/kubetop/ui/util"
	sview "github.com/boz/kubetop/ui/view/service"
)

func NewServiceTable(ctx elements.Context, ds service.Publisher) elements.Widget {
	ctx = ctx.New("service/table")
	content := table.NewWidget(ctx.Env(), sview.TableColumns(), true)

	ctx.AlsoClose(service.NewMonitor(ds,
		uiutil.ServicesPoster(ctx, sview.NewTable(content))))

	return elements.NewWidget(ctx, content)
}
