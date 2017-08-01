package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

func AlignRight(content views.Widget) views.Widget {
	return &aligner{content: content}
}

type aligner struct {
	content views.Widget
	view    views.View
}

func (w *aligner) Draw() {
	w.content.Draw()
}

func (w *aligner) Resize() {
	w.layout()
}

func (w *aligner) HandleEvent(ev tcell.Event) bool {
	return w.content.HandleEvent(ev)
}

func (w *aligner) SetView(view views.View) {
	w.view = view
	w.layout()
}

func (w *aligner) Size() (int, int) {
	return w.content.Size()
}

func (w *aligner) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *aligner) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}

func (w *aligner) layout() {
	if w.view == nil {
		return
	}

	vx, vy := w.view.Size()
	cx, _ := w.content.Size()

	xoff := 0
	if vx > cx {
		xoff = vx - cx
	}

	vp := views.NewViewPort(w.view, xoff, 0, vx-xoff, vy)
	w.content.SetView(vp)
	w.content.Resize()
}
