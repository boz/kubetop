package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Popupper struct {
	view    views.View
	current views.Widget
	views.WidgetWatchers
}

func NewPopupper() *Popupper {
	return &Popupper{}
}

// todo: stack

func (p *Popupper) Push(w views.Widget) {
	p.current = w
	p.current.SetView(p.view)
}

func (p *Popupper) Pop() views.Widget {
	current := p.current
	p.current = nil
	return current
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
	if p.current != nil && p.current.HandleEvent(ev) {
		return true
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEsc:
			p.Pop()
			return true
		}
	}
	return false
}
