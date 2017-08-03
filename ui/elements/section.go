package elements

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/util"
	"github.com/gdamore/tcell/views"
)

type Section interface {
	Themeable
}

func NewSection(env util.Env, title string, content views.Widget) Section {
	w := &section{
		title: views.NewText(),
		panes: panes{
			orientation: views.Vertical,
			expand:      false,
			env:         env,
		},
	}
	w.title.SetText(title)
	w.title.SetAlignment(views.HAlignCenter)
	w.Append(w.title)
	w.Append(content)
	return w
}

type section struct {
	panes
	theme theme.Theme
	title *views.Text
}

func (w *section) SetTheme(th theme.Theme) {
	w.title.SetStyle(th.Title)
	w.panes.SetTheme(th)
}

func (w *section) Theme() theme.Theme {
	return w.theme
}
