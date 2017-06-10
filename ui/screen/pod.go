package screen

import (
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/widget"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type podIndex struct {
	content elements.Widget
	ctx     elements.Context
}

func NewPodIndex(ctx elements.Context, ds pod.BaseDatasource) elements.Widget {
	ctx = ctx.New("pod/index")
	content := widget.NewPodTable(ctx, ds)
	index := &podIndex{content, ctx}
	content.Watch(index)
	return index
}

func (w *podIndex) Draw() {
	w.content.Draw()
}

func (w *podIndex) Resize() {
	w.content.Resize()
}

func (w *podIndex) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetResize:
		w.Resize()
		return true
	}
	return w.content.HandleEvent(ev)
}

func (w *podIndex) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *podIndex) Size() (int, int) {
	return w.content.Size()
}

func (w *podIndex) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *podIndex) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}

func (w *podIndex) Close() {
	w.ctx.Close()
}
