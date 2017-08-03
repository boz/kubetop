package deflist

import (
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const padsize = 1

type Widget interface {
	views.Widget
	theme.Themeable
	SetRows([]Row)
}

type Row interface {
	Term() views.Widget
	Definition() views.Widget
	SetTheme(theme.Theme)
}

type row struct {
	term       elements.Styleable
	definition elements.Styleable
	variant    theme.LabelVariant
	theme      theme.Theme
}

func NewSimpleRow(termtxt, deftxt string, lv theme.LabelVariant) Row {
	termw := views.NewText()
	termw.SetAlignment(views.HAlignLeft)
	termw.SetText(termtxt)

	defw := views.NewText()
	defw.SetAlignment(views.HAlignLeft)
	defw.SetText(deftxt)
	return NewRow(termw, defw, lv)
}

func NewRow(term, definition elements.Styleable, lv theme.LabelVariant) Row {
	return &row{term: term, definition: definition, variant: lv}
}

func (r *row) Term() views.Widget {
	return r.term
}

func (r *row) Definition() views.Widget {
	return r.definition
}

func (r *row) SetTheme(th theme.Theme) {
	r.term.SetStyle(th.Deflist.Term.Get(r.variant))
	r.definition.SetStyle(th.Deflist.Definition.Get(r.variant))
	r.theme = th
}

type widget struct {
	view   views.View
	rows   []Row
	width  int
	height int
	theme  theme.Theme
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

func (w *widget) SetTheme(th theme.Theme) {
	w.theme = th
	for _, row := range w.rows {
		row.SetTheme(th)
	}
}

func (w *widget) SetView(view views.View) {
	w.view = view
	w.layout()
}

func (w *widget) Size() (int, int) {
	return w.width, w.height
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

	w.width = termx + defx + padsize
	w.height = len(w.rows)
}
