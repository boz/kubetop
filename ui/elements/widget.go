package elements

import (
	"github.com/boz/kcache/nsname"
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Widget interface {
	views.Widget
	theme.Themeable
	Close()
}

type NSNameWidget interface {
	Widget
	ID() nsname.NSName
}

type widget struct {
	content views.Widget
	ctx     Context
	theme   theme.Theme
}

func NewWidget(ctx Context, content views.Widget) Widget {
	return &widget{content: content, ctx: ctx}
}

func (w *widget) SetTheme(th theme.Theme) {
	w.theme = th
	if content, ok := w.content.(theme.Themeable); ok {
		content.SetTheme(th)
	}
}

func (w *widget) Theme() theme.Theme {
	return w.theme
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

func NewNSNameWidget(ctx Context, content views.Widget, id nsname.NSName) NSNameWidget {
	return &nsNameWidget{
		widget: widget{content: content, ctx: ctx},
		id:     id,
	}
}

type nsNameWidget struct {
	widget
	id nsname.NSName
}

func (w *nsNameWidget) ID() nsname.NSName {
	return w.id
}
