package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Popupper struct {
	view    views.View
	current views.Widget
	views.WidgetWatchers

	Presentable
}

func NewPopupper(p Presenter) *Popupper {
	w := &Popupper{}
	p.New("ui/elements/popper", w)
	return w
}

// todo: stack

func (p *Popupper) Push(w views.Widget) {
	if w == nil {
		return
	}
	p.Pop()
	w.Watch(p)
	w.SetView(p.view)
	p.current = w
}

func (p *Popupper) Pop() views.Widget {
	prev := p.current
	if prev != nil {
		prev.Unwatch(p)
		p.current = nil
	}
	return prev
}

func (p *Popupper) Draw() {
	if p.current != nil {
		p.current.Draw()
	}
}

func (p *Popupper) Resize() {
	if p.current != nil {
		p.current.Resize()
	}
}

func (p *Popupper) Size() (int, int) {
	return 0, 0
}

func (p *Popupper) SetView(view views.View) {
	p.view = view
	if p.current != nil {
		p.current.SetView(view)
	}
}

func (p *Popupper) HandleEvent(ev tcell.Event) bool {

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

	return p.current.HandleEvent(ev)

}
