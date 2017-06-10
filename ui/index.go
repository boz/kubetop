package ui

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type IndexWidget interface {
	views.Widget
	ResetRows([]elements.TableRow)
	UpdateRow(elements.TableRow)
	RemoveRow(string)
}

type IndexProvider interface {
	Stop()
}

type IndexBuilder interface {
	Model() elements.Table
	Create(IndexWidget, <-chan struct{}) IndexProvider
}

type indexWidget struct {
	builder  IndexBuilder
	provider IndexProvider
	content  *elements.TableWidget
	elements.Presentable
}

func NewIndexWidget(
	name string, p elements.Presenter, builder IndexBuilder) views.Widget {
	index := &indexWidget{
		builder: builder,
		content: elements.NewTableWidget(builder.Model()),
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

func (w *indexWidget) ResetRows(rows []elements.TableRow) {
	w.PostFunc(func() {
		for _, row := range rows {
			w.content.AddRow(row)
		}
		w.Resize()
	})
}

func (w *indexWidget) UpdateRow(row elements.TableRow) {
	w.PostFunc(func() {
		w.content.AddRow(row)
		w.Resize()
	})
}

func (w *indexWidget) RemoveRow(id string) {
	w.PostFunc(func() {
		w.content.RemoveRow(id)
		w.Resize()
	})
}
