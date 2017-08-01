package pod

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/screen/requests"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type indexScreen struct {
	layout  elements.Panes
	table   elements.Widget
	summary elements.NSNameWidget
	ctx     elements.Context
	views.WidgetWatchers
}

func NewIndex(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("pod/index")

	pods, err := ctx.Backend().Pods()
	if err != nil {
		return nil, err
	}

	table := widget.NewPodTable(ctx, pods)

	layout := elements.NewVPanes(true)
	layout.PushBackWidget(table)

	index := &indexScreen{
		layout: layout,
		table:  table,
		ctx:    ctx,
	}

	layout.Watch(index)
	table.Watch(index)

	return elements.NewScreen(ctx, req, "Pods", index), nil
}

func (w *indexScreen) Draw() {
	w.layout.Draw()
}

func (w *indexScreen) Resize() {
	w.layout.Resize()
}

func (w *indexScreen) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *views.EventWidgetContent:
		w.PostEventWidgetContent(w)
		return true
	case *table.EventRowActive:
		w.showSummary(ev.Row().ID())
		return true
	case *table.EventRowInactive:
		w.removeSummary()
		return true
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEnter:
			if w.summary != nil {
				// navigate to pods/show
				w.ctx.NavigateTo(requests.PodShowRequest(w.summary.ID()))
				return true
			}
		}
	}
	return w.table.HandleEvent(ev)
}

func (w *indexScreen) SetView(view views.View) {
	w.layout.SetView(view)
}

func (w *indexScreen) Size() (int, int) {
	return w.layout.Size()
}

func (w *indexScreen) showSummary(id string) {
	w.removeSummary()
	summary, _ := widget.NewPodSummary(w.ctx, id)
	summary.Watch(w)
	w.summary = summary
	w.layout.PushFrontWidget(w.summary)
}

func (w *indexScreen) removeSummary() {
	if w.summary != nil {
		w.summary.Unwatch(w)
		w.layout.RemoveWidget(w.summary)
		w.summary.Close()
	}
}
