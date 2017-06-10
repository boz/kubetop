package table

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell/views"
)

var (
	styleColPad = 2
)

type View interface {
	Size() (int, int)
	SetText(string, theme.LabelVariant)
	SetContent(x, y int, ch rune, comb []rune, th theme.LabelVariant)
}

type _view struct {
	view   views.View
	xoff   int
	yoff   int
	width  int
	height int
	theme  theme.LabelTheme
}

func newView(view views.View, xoff, yoff, width, height int, theme theme.LabelTheme) View {
	return &_view{view, xoff, yoff, width, height, theme}
}

func (v *_view) Size() (int, int) {
	return v.width, v.height
}

func (v *_view) SetContent(x, y int, ch rune, comb []rune, th theme.LabelVariant) {
	v.view.SetContent(x+v.xoff, y+v.yoff, ch, comb, v.theme.Get(th))
}

func (v *_view) SetText(text string, th theme.LabelVariant) {
	for i, ch := range text {
		v.SetContent(i, 0, ch, nil, th)
	}
	for i := len(text); i < v.width; i++ {
		v.SetContent(i, 0, ' ', nil, th)
	}
}
