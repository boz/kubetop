package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/screen"
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/version"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type mainWidget struct {
	stopch chan<- bool
	panel  *views.Panel

	navbar *views.SimpleStyledTextBar

	content elements.Widget

	popupper *elements.Popupper

	navigator elements.Navigator

	ctx elements.Context
}

func newNavBar() *views.SimpleStyledTextBar {
	bar := views.NewSimpleStyledTextBar()
	bar.SetStyle(theme.AppHeader.Bar)
	bar.RegisterCenterStyle('N', theme.AppHeader.Bar)
	bar.RegisterCenterStyle('A', theme.AppHeader.Action)
	bar.SetCenter("%Nloading...")
	return bar
}

func newMainStatus() views.Widget {
	bar := views.NewSimpleStyledTextBar()
	bar.SetStyle(theme.AppHeader.Bar)

	bar.RegisterLeftStyle('N', theme.AppHeader.Bar)
	bar.RegisterLeftStyle('A', theme.AppHeader.Action)
	bar.SetLeft("%N[%AQ%N] Quit %N[%AP%N] Pods %N[%AS%N] Services %N[%AN%N] Nodes")

	bar.RegisterRightStyle('N', theme.AppHeader.Bar)
	bar.SetRight(fmt.Sprintf("%%Nkubetop %v", version.Version()))
	return bar
}

func newMainWidget(ctx elements.Context, stopch chan<- bool) views.Widget {
	ctx = ctx.New("ui/main")

	navbar := newNavBar()

	panel := views.NewPanel()
	panel.SetTitle(navbar)
	panel.SetStatus(newMainStatus())

	router := elements.NewRouter(ctx)
	screen.RegisterPodRoutes(router)
	screen.RegisterServiceRoutes(router)
	screen.RegisterNodeRoutes(router)
	screen.RegisterEventRoutes(router)

	widget := &mainWidget{
		stopch:    stopch,
		panel:     panel,
		navbar:    navbar,
		popupper:  elements.NewPopupper(ctx),
		ctx:       ctx,
		navigator: router,
	}

	ctx.WatchNavigation(widget)

	ctx.NavigateTo(screen.PodIndexRequest())

	return widget
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
				w.ctx.NavigateTo(screen.PodIndexRequest())
			case 'S', 's':
				w.ctx.NavigateTo(screen.ServiceIndexRequest())
			case 'N', 'n':
				w.ctx.NavigateTo(screen.NodeIndexRequest())
			case 'E', 'e':
				w.ctx.NavigateTo(screen.EventIndexRequest())
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

func (w *mainWidget) HandleNavigationRequest(req elements.Request) {
	screen, err := w.navigator.Open(req)
	if err != nil {
		w.ctx.Env().LogErr(err, "can't open request %v", req.Route())
		return
	}
	w.setContent(screen)
}

func (w *mainWidget) setContent(child elements.Screen) {
	if cur := w.content; cur != nil {
		cur.Close()
	}
	w.content = child
	w.panel.SetContent(child)
	w.navbar.SetCenter(child.State().Title())
	w.Resize()
}

func (w *mainWidget) textArea() views.Widget {
	var text string

	for i := 0; i < 9; i++ {
		text = text + strconv.Itoa(i) + " " + strings.Repeat("123456789", 2) + "\n"
	}
	text = text + text

	txt := views.NewTextArea()
	txt.SetContent("%N" + text)
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
