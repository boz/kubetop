package ui

import (
	"strconv"
	"strings"

	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/screen"
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type mainWidget struct {
	stopch chan<- bool
	panel  *views.Panel

	content elements.Widget

	popupper *elements.Popupper

	ctx elements.Context
}

func newMainTitle() views.Widget {
	title := views.NewSimpleStyledTextBar()
	title.SetStyle(theme.AppHeader.Bar)

	title.RegisterLeftStyle('N', theme.AppHeader.Bar)
	title.RegisterLeftStyle('A', theme.AppHeader.Action)
	title.SetLeft("%N[%AQ%N] Quit")

	title.RegisterRightStyle('N', theme.AppHeader.Bar)
	title.SetRight("%Nkubetop")
	return title
}

func newMainWidget(ctx elements.Context, stopch chan<- bool) views.Widget {
	panel := views.NewPanel()
	panel.SetTitle(newMainTitle())

	return &mainWidget{
		stopch:   stopch,
		panel:    panel,
		popupper: elements.NewPopupper(ctx),
		ctx:      ctx.New("ui/main"),
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
				w.showPodIndex()
			case 'S', 's':
				w.showServiceIndex()
			case 'X', 'x':
				popup := elements.NewPopup(w.ctx, 10, 10, theme.Base)
				popup.SetContent(w.textArea())
				w.popupper.Push(popup)
				return true
			}
		}
	}

	return false
}

func (w *mainWidget) showPodIndex() {
	ds, _ := w.ctx.Backend().Pods()
	widget := screen.NewPodIndex(w.ctx, ds)
	w.setContent(widget)
}

func (w *mainWidget) showServiceIndex() {
	ds, _ := w.ctx.Backend().Services()
	widget := screen.NewServiceIndex(w.ctx, ds)
	w.setContent(widget)
}

func (w *mainWidget) setContent(child elements.Widget) {
	if cur := w.content; cur != nil {
		cur.Close()
	}
	w.content = child
	w.panel.SetContent(child)
	w.Resize()
}

func (w *mainWidget) textArea() views.Widget {
	var text string

	for i := 0; i < 9; i++ {
		text = text + strconv.Itoa(i) + " " + strings.Repeat("123456789", 2) + "\n"
	}
	text = text + text

	txt := views.NewTextArea()
	txt.SetContent(text)
	txt.EnableCursor(true)
	return txt
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
