package widget

import (
	"github.com/boz/kubetop/backend/service"
	"github.com/boz/kubetop/ui/controller"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/view"
)

func NewServiceTable(ctx elements.Context, ds service.BaseDatasource) elements.Widget {
	content := table.NewWidget(view.ServiceTableColumns())
	ctx = ctx.New("service/table")
	handler := controller.NewServicesPostHandler(ctx, view.NewServiceTableWriter(content))
	controller.NewServiceController(ctx, ds, handler)
	return elements.NewWidget(ctx, content)
}
