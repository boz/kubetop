package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Popupper interface {
	views.Widget
	Push(w views.Widget)
	Pop() views.Widget
}

type popupper struct {
	view    views.View
	current views.Widget
	views.WidgetWatchers

	ctx Context
}

func NewPopupper(ctx Context) Popupper {
	w := &popupper{ctx: ctx.New("ui/elements/popupper")}
	return w
}

// todo: stack

func (p *popupper) Push(w views.Widget) {
	if w == nil {
		return
	}
	p.Pop()
	w.SetView(p.view)
	w.Watch(p)
	p.current = w
}

func (p *popupper) Pop() views.Widget {
	prev := p.current
	if prev != nil {
		prev.Unwatch(p)
		p.current = nil
	}
	return prev
}

func (p *popupper) Draw() {
	if p.current != nil {
		p.current.Draw()
	}
}

func (p *popupper) Resize() {
	if p.current != nil {
		p.current.Resize()
	}
}

func (p *popupper) Size() (int, int) {
	return 0, 0
}

func (p *popupper) SetView(view views.View) {
	p.view = view
	if p.current != nil {
		p.current.SetView(view)
	}
}

func (p *popupper) HandleEvent(ev tcell.Event) bool {

	if p.current == nil {
		return false
	}

	switch ev := ev.(type) {
	case *EventPopupClose:
		if ev.Widget() == p.current {
			p.Pop()
			return true
		}
	}

	if p.current.HandleEvent(ev) {
		return true
	}

	if _, ok := ev.(*tcell.EventKey); ok && p.current != nil {
		return true
	}

	return false
}
