package widget

import (
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/controller"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/view"
)

func NewPodTable(ctx elements.Context, ds pod.BaseDatasource) elements.Widget {
	content := table.NewWidget(view.PodTableColumns())
	ctx = ctx.New("pod/table")
	handler := controller.NewPodPostHandler(ctx, view.NewPodTableWriter(content))
	controller.NewPodController(ctx, ds, handler)
	return elements.NewWidget(ctx, content)
}
