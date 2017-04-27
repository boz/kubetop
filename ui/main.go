package ui

import (
	"strconv"
	"strings"

	"github.com/boz/kubetop/ui/elements"
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

	popupper *elements.Popupper
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
		stopch:   stopch,
		panel:    panel,
		popupper: elements.NewPopupper(),
	}
}

func (w *mainWidget) Draw() {
	w.panel.Draw()
	w.popupper.Draw()
}

func (w *mainWidget) Resize() {
	w.panel.Resize()
	w.popupper.Resize()
}

func (w *mainWidget) HandleEvent(ev tcell.Event) bool {

	if w.popupper.HandleEvent(ev) {
		return true
	}

	if w.panel.HandleEvent(ev) {
		return true
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'Q', 'q':
				w.stopch <- true
				return true
			case 'P', 'p':
				popup := elements.NewPopup(10, 10, tcell.StyleDefault)
				popup.SetContent(textArea())
				w.popupper.Push(popup)
				return true
			}
		}
	}

	return false
}

func textArea() views.Widget {
	text := ""

	for i := 0; i < 9; i++ {
		text = text + strconv.Itoa(i) + " " + strings.Repeat("123456789", 2) + "\n"
	}

	text = text + text

	w := views.NewTextArea()
	w.SetContent(text)
	w.EnableCursor(true)
	return w
}

func (w *mainWidget) SetView(view views.View) {
	w.panel.SetView(view)
	w.popupper.SetView(view)
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
