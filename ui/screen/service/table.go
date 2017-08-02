package service

import (
	"github.com/boz/kcache/types/service"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	sview "github.com/boz/kubetop/ui/view/service"
)

func newIndexTable(ctx elements.Context, ds service.Publisher) elements.Widget {
	ctx = ctx.New("service/index#table")
	content := table.NewWidget(ctx.Env(), sview.TableColumns(), true)
	sview.Monitor(ctx, ds, sview.NewTable(content))
	return elements.NewWidget(ctx, content)
}
