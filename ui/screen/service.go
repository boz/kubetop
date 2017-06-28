package screen

import (
	"github.com/boz/kubetop/backend/service"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type serviceIndex struct {
	content elements.Widget
	ctx     elements.Context
}

func NewServiceIndex(ctx elements.Context, ds service.BaseDatasource) elements.Widget {
	ctx = ctx.New("service/index")
	content := widget.NewServiceTable(ctx, ds)
	index := &serviceIndex{content, ctx}
	content.Watch(index)
	return index
}

func (w *serviceIndex) Draw() {
	w.content.Draw()
}

func (w *serviceIndex) Resize() {
	w.content.Resize()
}

func (w *serviceIndex) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		w.Resize()
		return true
	}
	return w.content.HandleEvent(ev)
}

func (w *serviceIndex) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *serviceIndex) Size() (int, int) {
	return w.content.Size()
}

func (w *serviceIndex) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *serviceIndex) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}

func (w *serviceIndex) Close() {
	w.ctx.Close()
}
