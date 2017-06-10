package table

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell/views"
)

var (
	styleColPad = 2
)

type CellView interface {
	Size() (int, int)
	SetText(string, theme.LabelVariant)
	SetContent(x, y int, ch rune, comb []rune, th theme.LabelVariant)
}

type cellView struct {
	view   views.View
	xoff   int
	yoff   int
	width  int
	height int
	theme  theme.LabelTheme
}

func NewCellView(view views.View, xoff, yoff, width, height int, theme theme.LabelTheme) CellView {
	return &cellView{view, xoff, yoff, width, height, theme}
}

func (v *cellView) Size() (int, int) {
	return v.width, v.height
}

func (v *cellView) SetContent(x, y int, ch rune, comb []rune, th theme.LabelVariant) {
	v.view.SetContent(x+v.xoff, y+v.yoff, ch, comb, v.theme.Get(th))
}

func (v *cellView) SetText(text string, th theme.LabelVariant) {
	for i, ch := range text {
		v.SetContent(i, 0, ch, nil, th)
	}
	for i := len(text); i < v.width; i++ {
		v.SetContent(i, 0, ' ', nil, th)
	}
}
