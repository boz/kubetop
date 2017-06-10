package table

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	styleTH     = tcell.StyleDefault.Bold(true)
	styleColPad = 2
)

type CellView interface {
	Size() (int, int)
	SetText(string, tcell.Style)
	SetContent(x, y int, ch rune, comb []rune, s tcell.Style)
}

type cellView struct {
	view     views.View
	xoff     int
	yoff     int
	width    int
	height   int
	selected bool
}

func NewCellView(view views.View, xoff, yoff, width, height int, selected bool) CellView {
	return &cellView{view, xoff, yoff, width, height, selected}
}

func (v *cellView) Size() (int, int) {
	return v.width, v.height
}

func (v *cellView) SetContent(x, y int, ch rune, comb []rune, s tcell.Style) {
	v.view.SetContent(x+v.xoff, y+v.yoff, ch, comb, s.Reverse(v.selected))
}

func (v *cellView) SetText(text string, s tcell.Style) {
	for i, ch := range text {
		v.SetContent(i, 0, ch, nil, s)
	}
	for i := len(text); i < v.width; i++ {
		v.SetContent(i, 0, ' ', nil, s)
	}
}
