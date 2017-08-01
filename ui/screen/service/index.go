package service

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type indexScreen struct {
	content elements.Widget
	ctx     elements.Context
}

func NewIndex(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("service/index")

	db, err := ctx.Backend().Services()
	if err != nil {
		return nil, err
	}
	content := widget.NewServiceTable(ctx, db)
	index := &indexScreen{content, ctx}
	content.Watch(index)

	return elements.NewScreen(ctx, req, "Services", index), nil
}

func (w *indexScreen) Draw() {
	w.content.Draw()
}

func (w *indexScreen) Resize() {
	w.content.Resize()
}

func (w *indexScreen) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		w.Resize()
		return true
	}
	return w.content.HandleEvent(ev)
}

func (w *indexScreen) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *indexScreen) Size() (int, int) {
	return w.content.Size()
}

func (w *indexScreen) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *indexScreen) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}
