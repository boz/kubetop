package pod

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/help"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/screen/requests"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type indexScreen struct {
	layout  elements.Sections
	summary elements.NSNameWidget
	ctx     elements.Context
	views.WidgetWatchers
}

func NewIndex(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("pod/index")

	ds, err := newIndexDS(ctx)
	if err != nil {
		return nil, err
	}

	table := newIndexTable(ctx, ds.pods)

	layout := elements.NewVSections(ctx.Env(), true)
	layout.Append(table)

	index := &indexScreen{
		layout: layout,
		ctx:    ctx,
	}

	layout.Watch(index)
	table.Watch(index)

	hsections := []views.Widget{
		help.NewSection(ctx.Env(), "Pod List", []help.Key{
			help.NewKey("ctrl-k", "Kill pod"),
			help.NewKey("ctrl-l", "View logs"),
			help.NewKey("enter", "View details"),
		}),
	}

	return elements.NewScreen(ctx, req, "Pods", index, hsections), nil
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
	return w.layout.HandleEvent(ev)
}

func (w *indexScreen) SetView(view views.View) {
	w.layout.SetView(view)
}

func (w *indexScreen) Size() (int, int) {
	return w.layout.Size()
}

func (w *indexScreen) showSummary(id string) {
	w.removeSummary()
	summary, err := newSummary(w.ctx, id)
	if err != nil {
		w.ctx.Env().LogErr(err, "error opening summary")
		return
	}
	summary.Watch(w)
	w.summary = summary
	w.layout.Prepend(w.summary)
}

func (w *indexScreen) removeSummary() {
	if w.summary != nil {
		w.summary.Unwatch(w)
		w.layout.Remove(w.summary)
		w.summary.Close()
	}
}
