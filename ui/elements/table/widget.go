package table

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Widget struct {
	model *tableModel
	colsz []int

	view  views.View
	hport *views.ViewPort
	rport *views.ViewPort
	views.WidgetWatchers
}

func NewWidget(cols []TH) *Widget {
	return &Widget{
		model: newTableModel(cols),
		hport: views.NewViewPort(nil, 0, 0, 0, 0),
		rport: views.NewViewPort(nil, 0, 1, 0, 0),
	}
}

func (tw *Widget) ResetRows(rows []TR) {
	tw.model.reset(rows)
	tw.resizeContent()
}

func (tw *Widget) InsertRow(row TR) {
	tw.model.insert(row)
	tw.resizeContent()
}

func (tw *Widget) UpdateRow(row TR) {
	tw.model.update(row)
	tw.resizeContent()
}

func (tw *Widget) RemoveRow(id string) {
	tw.model.remove(id)
	tw.resizeContent()
}

func (tw *Widget) Draw() {
	tw.hport.Fill(' ', theme.Base)
	tw.rport.Fill(' ', theme.Base)
	tw.drawHeader()
	tw.model.each(func(roff int, row TR) {
		tw.drawRow(roff, row)
	})
}

func (tw *Widget) Resize() {
	if tw.view == nil {
		return
	}
	tw.resizeContent()
}

func (tw *Widget) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {

		case tcell.KeyUp, tcell.KeyCtrlP:
			return tw.keyUp()
		case tcell.KeyDown, tcell.KeyCtrlN:
			return tw.keyDown()
		case tcell.KeyLeft, tcell.KeyCtrlB:
			return tw.keyLeft()
		case tcell.KeyRight, tcell.KeyCtrlF:
			return tw.keyRight()
		case tcell.KeyEscape:
			return tw.keyEscape()

		case tcell.KeyRune:
			switch ev.Rune() {
			case 'k':
				return tw.keyUp()
			case 'j':
				return tw.keyDown()
			case 'h':
				return tw.keyLeft()
			case 'l':
				return tw.keyRight()
			}
		}
	}
	return false
}

func (tw *Widget) SetView(view views.View) {
	tw.view = view
	tw.hport.SetView(view)
	tw.rport.SetView(view)
	tw.Resize()
}

func (tw *Widget) Size() (int, int) {
	x, y := tw.rport.Size()
	return x, y + 1
}

func (tw *Widget) resizeContent() {
	colsz := make([]int, len(tw.model.columns()))

	update := func(i int, col TD) {
		width, _ := col.Size()
		if width+styleColPad > colsz[i] {
			colsz[i] = width + styleColPad
		}
	}

	for i, col := range tw.model.columns() {
		update(i, col)
	}

	tw.model.each(func(_ int, row TR) {
		for i, col := range row.Columns() {
			update(i, col)
		}
	})

	width := 0
	for _, col := range colsz {
		width += col
	}
	height := len(tw.model.columns())

	tw.hport.Resize(0, 0, width, 1)
	tw.rport.Resize(0, 1, width, height)

	tw.colsz = colsz
}

func (tw *Widget) drawHeader() {
	xoff := 0
	yoff := 0
	cols := tw.model.columns()
	view := tw.hport
	for i, col := range cols {
		width := tw.colsz[i]
		cview := NewCellView(view, xoff, yoff, width, 1, theme.Table.TH)
		col.Draw(cview)
		xoff += width
	}
}

func (tw *Widget) drawRow(yoff int, row TR) {
	xoff := 0
	cols := row.Columns()
	view := tw.rport

	lth := theme.Table.TD
	if tw.model.isSelected(row.ID()) {
		lth = theme.Table.TDSelected
	}

	for i, col := range cols {
		width := tw.colsz[i]
		cview := NewCellView(view, xoff, yoff, width, 1, lth)
		col.Draw(cview)
		xoff += width
	}
}

func (tw *Widget) keyUp() bool {
	idx, _ := tw.model.selectPrev()
	if idx < 0 {
		return false
	}
	tw.rport.MakeVisible(-1, idx)
	return true
}

func (tw *Widget) keyDown() bool {
	idx, _ := tw.model.selectNext()
	if idx < 0 {
		return false
	}
	tw.rport.MakeVisible(-1, idx)
	return true
}

func (tw *Widget) keyLeft() bool {
	tw.hport.ScrollLeft(1)
	tw.rport.ScrollLeft(1)
	return true
}
func (tw *Widget) keyRight() bool {
	tw.hport.ScrollRight(1)
	tw.rport.ScrollRight(1)
	return true
}

func (tw *Widget) keyEscape() bool {
	return tw.model.clearSelection()
}
