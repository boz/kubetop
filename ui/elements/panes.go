package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Panes interface {
	views.Widget

	Widgets() []views.Widget

	PushBackWidget(views.Widget)
	PushFrontWidget(views.Widget)
	RemoveWidget(views.Widget)
	InsertBeforeWidget(views.Widget, views.Widget)
	InsertAfterWidget(views.Widget, views.Widget)
}

type panesChild struct {
	widget views.Widget
	view   *views.ViewPort
}

type panes struct {
	children []*panesChild
	view     views.View

	views.WidgetWatchers
}

func NewPanes() Panes {
	return &panes{}
}

func (p *panes) Widgets() []views.Widget {
	children := make([]views.Widget, 0, len(p.children))
	for _, c := range p.children {
		children = append(children, c.widget)
	}
	return children
}

func (p *panes) Draw() {
	for _, c := range p.children {
		c.widget.Draw()
	}
}

func (p *panes) Resize() {
	p.layout()

	for _, c := range p.children {
		c.widget.Resize()
	}

	p.PostEventWidgetResize(p)
}

func (p *panes) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		p.layout()
		p.PostEventWidgetContent(p)
		return true
	}
	for _, c := range p.children {
		if c.widget.HandleEvent(ev) {
			return true
		}
	}
	return false
}

func (p *panes) SetView(view views.View) {
	p.view = view
	for _, c := range p.children {
		c.view.SetView(view)
	}
}

func (p *panes) Size() (int, int) {
	px, py := 0, 0

	for _, c := range p.children {
		cx, cy := c.widget.Size()

		py += cy
		if cx > px {
			px = cx
		}
	}

	return px, py
}

func (p *panes) PushBackWidget(w views.Widget) {
	cnew := p.newChild(w)
	p.children = append(p.children, cnew)
	p.afterModify()
}

func (p *panes) PushFrontWidget(w views.Widget) {
	cnew := p.newChild(w)
	p.children = append([]*panesChild{cnew}, p.children...)
	p.afterModify()
}

func (p *panes) RemoveWidget(w views.Widget) {
	changed := false

	for i, c := range p.children {
		if c.widget == w {
			changed = true
			p.children = append(p.children[:i], p.children[i+1:]...)
		}
	}

	if changed {
		p.afterModify()
	}
}

func (p *panes) InsertBeforeWidget(mark views.Widget, w views.Widget) {
	for i, c := range p.children {
		if c.widget == mark {
			cnew := p.newChild(w)

			p.children = append(p.children, nil)
			copy(p.children[i+1:], p.children[i:])
			p.children[i] = cnew

			p.afterModify()
			return
		}
	}
}

func (p *panes) InsertAfterWidget(mark views.Widget, w views.Widget) {
	for i, c := range p.children {
		if c.widget == mark {
			cnew := p.newChild(w)

			p.children = append(p.children, nil)
			copy(p.children[i+2:], p.children[i+1:])
			p.children[i+1] = cnew

			p.afterModify()
			return
		}
	}
}

func (p *panes) newChild(w views.Widget) *panesChild {
	c := &panesChild{
		view:   views.NewViewPort(p.view, 0, 0, 0, 0),
		widget: w,
	}
	w.SetView(c.view)
	w.Watch(p)
	return c
}

func (p *panes) afterModify() {
	p.layout()
	p.PostEventWidgetContent(p)
}

func (p *panes) layout() {
	if p.view == nil {
		return
	}

	vx, vy := p.view.Size()

	py := 0

	for i, c := range p.children {
		_, wy := c.widget.Size()

		if wy+py > vy {
			wy = vy - py
		}

		if i == len(p.children)-1 {
			wy = vy - py
		}

		c.view.Resize(0, py, vx, wy)
		c.widget.Resize()

		py += wy
	}
}
