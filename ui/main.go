package ui

import (
	"fmt"

	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type mainWidget struct {
	stopch chan<- bool
	panel  *views.Panel

	content views.Widget

	popupper *elements.Popupper

	elements.Presentable
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

func newMainWidget(p elements.Presenter, stopch chan<- bool) views.Widget {

	panel := views.NewPanel()
	panel.SetTitle(newMainTitle())

	widget := &mainWidget{
		stopch:   stopch,
		panel:    panel,
		popupper: elements.NewPopupper(p),
	}

	return p.New("ui/main", widget)
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
			case 'X', 'x':
				popup := elements.NewPopup(w.Presenter(), 10, 10, theme.Base)
				popup.SetContent(w.textArea())
				w.popupper.Push(popup)
				return true
			}
		}
	}

	return false
}

func (w *mainWidget) showPodIndex() {
	ds, err := w.Backend().Pods()
	if err != nil {
	}
	env := w.Env().ForComponent("pods/index")
	builder := newPodIndexBuilder(env, ds)
	widget := NewIndexWidget("pods/index", w.Presenter(), builder)
	w.setContent(widget)
}

func (w *mainWidget) setContent(child views.Widget) {
	if cur := w.content; cur != nil {
		elements.ClosePresenter(cur)
	}
	w.content = child
	w.panel.SetContent(child)
}

func (w *mainWidget) textArea() views.Widget {
	var pods []pod.Pod
	var err error
	var text string

	src, err := w.Presenter().Backend().Pods()

	if err != nil {
		w.Env().LogErr(err, "getting datasource")
		text += fmt.Sprint("ERROR", err)
		goto done
	}

	pods, err = src.List()
	if err != nil {
		w.Env().LogErr(err, "getting list")
		text += fmt.Sprint("ERROR", err)
		goto done
	}

	for _, pod := range pods {
		text += pod.Resource().GetName() + "\n"
	}

	/*
		for i := 0; i < 9; i++ {
			text = text + strconv.Itoa(i) + " " + strings.Repeat("123456789", 2) + "\n"
		}
		text = text + text
	*/

done:

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
