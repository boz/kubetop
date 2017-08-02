package ui

import (
	"fmt"

	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/elements/deflist"
	"github.com/boz/kubetop/ui/screen"
	"github.com/boz/kubetop/ui/screen/requests"
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

	popupper elements.Popupper

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

var kbdNav = []struct {
	key   string
	label string
}{
	{"?", "Help"},
	{"P", "Pods"},
	{"S", "Services"},
	{"N", "Nodes"},
	{"E", "Events"},
	{"Q", "Quit"},
}

func newMainStatus() views.Widget {
	bar := views.NewSimpleStyledTextBar()
	bar.SetStyle(theme.AppHeader.Bar)

	bar.RegisterLeftStyle('N', theme.AppHeader.Bar)
	bar.RegisterLeftStyle('A', theme.AppHeader.Action)

	leftNav := ""
	for _, nav := range kbdNav {
		leftNav += fmt.Sprintf(" %%N[%%A%v%%N] %v", nav.key, nav.label)
	}
	bar.SetLeft(leftNav)

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
	screen.RegisterRoutes(router)

	widget := &mainWidget{
		stopch:    stopch,
		panel:     panel,
		navbar:    navbar,
		popupper:  elements.NewPopupper(ctx),
		ctx:       ctx,
		navigator: router,
	}

	ctx.WatchNavigation(widget)

	ctx.NavigateTo(requests.PodIndexRequest())

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
			case '?':
				w.openHelp()
				return true
			case 'P':
				w.ctx.NavigateTo(requests.PodIndexRequest())
				return true
			case 'S':
				w.ctx.NavigateTo(requests.ServiceIndexRequest())
				return true
			case 'N':
				w.ctx.NavigateTo(requests.NodeIndexRequest())
				return true
			case 'E':
				w.ctx.NavigateTo(requests.EventIndexRequest())
				return true
			case 'T':
				w.ctx.NavigateTo(requests.ThemeIndexRequest())
				return true
			case 'Q':
				w.stopch <- true
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

func (w *mainWidget) openHelp() {
	popup := elements.NewPopup(w.ctx, theme.Popup.Normal, w.helpWidget())
	w.popupper.Push(popup)
}

func (w *mainWidget) helpWidget() views.Widget {
	rows := make([]deflist.Row, 0, len(kbdNav))

	for _, row := range kbdNav {
		roww := deflist.NewSimpleRow(row.key, row.label, theme.ThemeActive.Deflist)
		rows = append(rows, roww)
	}
	return deflist.NewWidget(rows)
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
