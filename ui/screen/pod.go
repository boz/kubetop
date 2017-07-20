package screen

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const (
	podIndexPath = "/pods"
	podShowPath  = "/pods/show"
)

func RegisterPodRoutes(router elements.Router) {
	router.Register(elements.NewRoute(podIndexPath), elements.NewHandler(podIndexHandler))
}

func PodIndexRequest() elements.Request {
	return elements.NewRequest(podIndexPath)
}

type podIndex struct {
	layout  *views.BoxLayout
	table   elements.Widget
	details elements.Widget
	ctx     elements.Context
}

func podIndexHandler(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("pod/index")

	pods, err := ctx.Backend().Pods()
	if err != nil {
		return nil, err
	}

	table := widget.NewPodTable(ctx, pods)

	layout := views.NewBoxLayout(views.Vertical)
	layout.AddWidget(table, 0.3)
	layout.AddWidget(views.NewSpacer(), 0.0)

	index := &podIndex{
		layout: layout,
		table:  table,
		ctx:    ctx,
	}

	layout.Watch(index)
	table.Watch(index)

	return elements.NewScreen(ctx, req, "Pods", index), nil
}

func (w *podIndex) Draw() {
	w.layout.Draw()
}

func (w *podIndex) Resize() {
	w.layout.Resize()
}

func (w *podIndex) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *views.EventWidgetContent:
		if ev.Widget() == w.layout {
			w.Resize()
			return true
		}
	case *table.EventRowActive:
		w.showDetails(ev.Row().ID())
		return true
	case *table.EventRowInactive:
		w.removeDetails()
		return true
	}
	return w.table.HandleEvent(ev)
}

func (w *podIndex) SetView(view views.View) {
	w.layout.SetView(view)
}

func (w *podIndex) Size() (int, int) {
	return w.layout.Size()
}

func (w *podIndex) Watch(handler tcell.EventHandler) {
	w.layout.Watch(handler)
}

func (w *podIndex) Unwatch(handler tcell.EventHandler) {
	w.layout.Unwatch(handler)
}

func (w *podIndex) showDetails(id string) {
	w.removeDetails()
	details, _ := widget.NewPodDetails(w.ctx, id)
	details.Watch(w)
	w.details = details
	w.layout.InsertWidget(0, w.details, 0.2)
}

func (w *podIndex) removeDetails() {
	if w.details != nil {
		w.details.Unwatch(w)
		w.layout.RemoveWidget(w.details)
		w.details.Close()
	}
}
