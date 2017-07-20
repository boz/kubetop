package screen

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const (
	serviceIndexPath = "/service"
	serviceShowPath  = "/service/show"
)

func RegisterServiceRoutes(router elements.Router) {
	router.Register(elements.NewRoute(serviceIndexPath), elements.NewHandler(serviceIndexHandler))
}

func ServiceIndexRequest() elements.Request {
	return elements.NewRequest(serviceIndexPath)
}

type serviceIndex struct {
	content elements.Widget
	ctx     elements.Context
}

func serviceIndexHandler(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("service/index")

	db, err := ctx.Backend().Services()
	if err != nil {
		return nil, err
	}
	content := widget.NewServiceTable(ctx, db)
	index := &serviceIndex{content, ctx}
	content.Watch(index)

	return elements.NewScreen(ctx, req, "Services", index), nil
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
