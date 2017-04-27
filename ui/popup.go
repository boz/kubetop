package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type popup struct {
	view views.View

	content views.Widget
	cview   *views.ViewPort

	style tcell.Style

	width  int
	height int
}

func newPopup(width int, height int, style tcell.Style) *popup {
	return &popup{
		width:  width,
		height: height,
		style:  style,
		cview:  views.NewViewPort(nil, 1, 1, width-2, height-2),
	}
}

func (p *popup) Draw() {
}

func (p *popup) Resize() {
}

func (p *popup) HandleEvent(ev tcell.Event) bool {
	return false
}

func (p *popup) SetView(view views.View) {
	p.view = view
	p.cview.SetView(view)
}

func (p *popup) Size() (int, int) {
	return p.width, p.height
}

func (p *popup) Watch(handler tcell.EventHandler) {
}

func (w *popup) Unwatch(handler tcell.EventHandler) {
}

func (p *popup) SetContent(w views.Widget) {
	p.content = w
	w.SetView(p.cview)
}
