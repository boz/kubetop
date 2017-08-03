package event

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type indexScreen struct {
	layout elements.Sections
	ctx    elements.Context
}

func NewIndex(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("event/index")

	db, err := ctx.Backend().Events()
	if err != nil {
		return nil, err
	}

	layout := elements.NewVSections(ctx.Env(), true)
	table := newIndexTable(ctx, db)
	layout.Append(table)

	index := &indexScreen{layout, ctx}

	layout.Watch(index)
	table.Watch(index)

	return elements.NewScreen(ctx, req, "Events", index), nil
}

func (w *indexScreen) Draw() {
	w.layout.Draw()
}

func (w *indexScreen) Resize() {
	w.layout.Resize()
}

func (w *indexScreen) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		w.Resize()
		return true
	}
	return w.layout.HandleEvent(ev)
}

func (w *indexScreen) SetView(view views.View) {
	w.layout.SetView(view)
}

func (w *indexScreen) Size() (int, int) {
	return w.layout.Size()
}

func (w *indexScreen) Watch(handler tcell.EventHandler) {
	w.layout.Watch(handler)
}

func (w *indexScreen) Unwatch(handler tcell.EventHandler) {
	w.layout.Unwatch(handler)
}
