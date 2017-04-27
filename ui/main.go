package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	styleHeader = tcell.StyleDefault.
			Background(tcell.ColorTeal).
			Foreground(tcell.ColorGray)
	styleHeaderAction = tcell.StyleDefault.
				Background(tcell.ColorTeal).
				Foreground(tcell.ColorRed)
)

type mainWidget struct {
	stopch chan<- bool
	panel  *views.Panel
}

func newMainTitle() views.Widget {
	title := views.NewSimpleStyledTextBar()
	title.SetStyle(styleHeader)

	title.RegisterLeftStyle('N', styleHeader)
	title.RegisterLeftStyle('A', styleHeaderAction)
	title.SetLeft("%N[%AQ%N] Quit")

	title.RegisterRightStyle('N', styleHeader)
	title.SetRight("%Nkubetop")
	return title
}

func newMainWidget(stopch chan<- bool) *mainWidget {
	panel := views.NewPanel()
	panel.SetTitle(newMainTitle())
	return &mainWidget{
		stopch: stopch,
		panel:  panel,
	}
}

func (w *mainWidget) Draw() {
	w.panel.Draw()
}
func (w *mainWidget) Resize() {
	w.panel.Resize()
}

func (w *mainWidget) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'Q', 'q':
				w.stopch <- true
				return true
			}
		}
	}
	return w.panel.HandleEvent(ev)
}

func (w *mainWidget) SetView(view views.View) {
	w.panel.SetView(view)
}
func (w *mainWidget) Size() (int, int) {
	return w.panel.Size()
}
func (w *mainWidget) Watch(handler tcell.EventHandler) {
	w.panel.Watch(handler)
}
func (w *mainWidget) Unwatch(handler tcell.EventHandler) {
	w.panel.Unwatch(handler)
}
