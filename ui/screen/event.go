package screen

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const (
	eventIndexPath = "/event"
	eventShowPath  = "/event/show"
)

func RegisterEventRoutes(router elements.Router) {
	router.Register(elements.NewRoute(eventIndexPath), elements.NewHandler(eventIndexHandler))
}

func EventIndexRequest() elements.Request {
	return elements.NewRequest(eventIndexPath)
}

type eventIndex struct {
	content elements.Widget
	ctx     elements.Context
}

func eventIndexHandler(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("event/index")

	db, err := ctx.Backend().Events()
	if err != nil {
		return nil, err
	}

	content := widget.NewEventTable(ctx, db)
	index := &eventIndex{content, ctx}
	content.Watch(index)

	return elements.NewScreen(ctx, req, "Events", index), nil
}

func (w *eventIndex) Draw() {
	w.content.Draw()
}

func (w *eventIndex) Resize() {
	w.content.Resize()
}

func (w *eventIndex) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		w.Resize()
		return true
	}
	return w.content.HandleEvent(ev)
}

func (w *eventIndex) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *eventIndex) Size() (int, int) {
	return w.content.Size()
}

func (w *eventIndex) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *eventIndex) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}
