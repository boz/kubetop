package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Widget interface {
	views.Widget
	Close()
}

type widget struct {
	content views.Widget
	ctx     Context
}

func NewWidget(ctx Context, content views.Widget) Widget {
	return &widget{content, ctx}
}

func (w *widget) Draw() {
	w.content.Draw()
}

func (w *widget) Resize() {
	w.content.Resize()
}

func (w *widget) HandleEvent(ev tcell.Event) bool {
	return w.content.HandleEvent(ev)
}

func (w *widget) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *widget) Size() (int, int) {
	return w.content.Size()
}

func (w *widget) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *widget) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}

func (w *widget) Close() {
	w.ctx.Close()
}
