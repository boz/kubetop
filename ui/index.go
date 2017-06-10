package ui

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type IndexWidget interface {
	views.Widget
	ResetRows([]table.TR)
	InsertRow(table.TR)
	UpdateRow(table.TR)
	RemoveRow(string)
}

type IndexProvider interface {
	Stop()
}

type IndexBuilder interface {
	Model() []table.TH
	Create(IndexWidget, <-chan struct{}) IndexProvider
}

type indexWidget struct {
	builder  IndexBuilder
	provider IndexProvider
	content  *table.Widget
	elements.Presentable
}

func NewIndexWidget(
	name string, p elements.Presenter, builder IndexBuilder) views.Widget {
	index := &indexWidget{
		builder: builder,
		content: table.NewWidget(builder.Model()),
	}
	p.New(name, index)
	index.provider = builder.Create(index, index.Presenter().Closed())
	return index
}

func (w *indexWidget) Draw() {
	w.content.Draw()
}

func (w *indexWidget) Resize() {
	w.content.Resize()
}

func (w *indexWidget) HandleEvent(ev tcell.Event) bool {
	return w.content.HandleEvent(ev)
}

func (w *indexWidget) SetView(view views.View) {
	w.content.SetView(view)
}

func (w *indexWidget) Size() (int, int) {
	return w.content.Size()
}

func (w *indexWidget) Watch(handler tcell.EventHandler) {
	w.content.Watch(handler)
}

func (w *indexWidget) Unwatch(handler tcell.EventHandler) {
	w.content.Unwatch(handler)
}

func (w *indexWidget) ResetRows(rows []table.TR) {
	w.PostFunc(func() {
		w.content.ResetRows(rows)
		w.Resize()
	})
}

func (w *indexWidget) InsertRow(row table.TR) {
	w.PostFunc(func() {
		w.content.InsertRow(row)
		w.Resize()
	})
}

func (w *indexWidget) UpdateRow(row table.TR) {
	w.PostFunc(func() {
		w.content.UpdateRow(row)
		w.Resize()
	})
}

func (w *indexWidget) RemoveRow(id string) {
	w.PostFunc(func() {
		w.content.RemoveRow(id)
		w.Resize()
	})
}
