package deflist

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const padsize = 1

type Widget interface {
	views.Widget

	SetRows([]Row)
}

type Row interface {
	Term() views.Widget
	Definition() views.Widget
}

type row struct {
	term       views.Widget
	definition views.Widget
}

func NewSimpleRow(termtxt, deftxt string) Row {
	termw := views.NewText()
	termw.SetAlignment(views.HAlignLeft)
	termw.SetText(termtxt)

	defw := views.NewText()
	defw.SetAlignment(views.HAlignLeft)
	defw.SetText(deftxt)

	return NewRow(termw, defw)
}

func NewRow(term, definition views.Widget) Row {
	return &row{term, definition}
}

func (r *row) Term() views.Widget {
	return r.term
}

func (r *row) Definition() views.Widget {
	return r.definition
}

type widget struct {
	view views.View
	rows []Row
	views.WidgetWatchers
}

func NewWidget(rows []Row) Widget {
	w := &widget{}
	w.SetRows(rows)
	return w
}

func (w *widget) Draw() {
	for _, row := range w.rows {
		row.Term().Draw()
		row.Definition().Draw()
	}
}

func (w *widget) Resize() {
	w.layout()
	w.PostEventWidgetResize(w)
}

func (w *widget) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		w.layout()
		w.PostEventWidgetContent(w)
		return true
	}

	for _, row := range w.rows {
		switch {
		case row.Term().HandleEvent(ev):
			return true
		case row.Definition().HandleEvent(ev):
			return true
		}
	}
	return false
}

func (w *widget) SetView(view views.View) {
	w.view = view
	w.layout()
}

func (w *widget) Size() (int, int) {

	wx := 0

	for _, row := range w.rows {
		tx, _ := row.Term().Size()
		dx, _ := row.Definition().Size()
		wx += tx + dx + padsize
	}

	return wx, len(w.rows)
}

func (w *widget) SetRows(rows []Row) {
	for _, row := range w.rows {
		row.Term().Unwatch(w)
		row.Definition().Unwatch(w)
	}

	for _, row := range rows {
		row.Term().Unwatch(w)
		row.Definition().Unwatch(w)
	}
	w.rows = rows

	w.layout()
	w.PostEventWidgetContent(w)
}

func (w *widget) layout() {

	if w.view == nil {
		return
	}

	termx := 0
	defx := 0

	for _, row := range w.rows {
		tx, _ := row.Term().Size()
		if termx < tx {
			termx = tx
		}

		dx, _ := row.Definition().Size()
		if defx < dx {
			defx = dx
		}
	}

	for i, row := range w.rows {
		termv := views.NewViewPort(w.view, 0, i, termx, 1)
		defv := views.NewViewPort(w.view, termx+padsize, i, defx, 1)

		row.Term().SetView(termv)
		row.Term().Resize()

		row.Definition().SetView(defv)
		row.Definition().Resize()
	}
}
