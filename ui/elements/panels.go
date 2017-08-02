package elements

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type SelectableWidget interface {
	views.Widget
	theme.Themeable
}

type Panels interface {
	views.Widget

	Widgets() []SelectableWidget

	Append(SelectableWidget)
	Prepend(SelectableWidget)
	Remove(SelectableWidget)
	InsertBefore(SelectableWidget, SelectableWidget)
	InsertAfter(SelectableWidget, SelectableWidget)

	Selected() SelectableWidget
}

type panels struct {
	panes
	selected SelectableWidget
}

func NewVPanels(expand bool) Panels {
	return NewPanels(views.Vertical, expand)
}

func NewHPanels(expand bool) Panels {
	return NewPanels(views.Horizontal, expand)
}

func NewPanels(o views.Orientation, expand bool) Panels {
	return &panels{
		panes: panes{orientation: o, expand: expand},
	}
}

func (p *panels) Selected() SelectableWidget {
	return p.selected
}

func (p *panels) Widgets() []SelectableWidget {
	children := make([]SelectableWidget, 0, len(p.children))
	for _, c := range p.children {
		children = append(children, c.widget.(SelectableWidget))
	}
	return children
}

func (p *panels) Append(w SelectableWidget) {
	p.panes.Append(w)
	p.afterAdd(w)
}

func (p *panels) Prepend(w SelectableWidget) {
	p.panes.Prepend(w)
	p.afterAdd(w)
}

func (p *panels) Remove(w SelectableWidget) {
	if p.selected == nil || p.selected != w {
		p.panes.Remove(w)
		return
	}
	for idx, c := range p.children {
		if p.selected != c.widget.(SelectableWidget) {
			continue
		}
		if idx >= len(p.children)-1 {
			continue
		}

		wnext := p.children[idx+1].widget.(SelectableWidget)

		if wnext == p.selected {
			continue
		}

		p.selected = wnext
		p.selected.SetTheme(theme.ThemeActive)
		break
	}
	p.panes.Remove(w)
}

func (p *panels) InsertBefore(mark SelectableWidget, w SelectableWidget) {
	p.panes.InsertBefore(mark, w)
}

func (p *panels) InsertAfter(mark SelectableWidget, w SelectableWidget) {
	p.panes.InsertAfter(mark, w)
}

func (p *panels) HandleEvent(ev tcell.Event) bool {

	switch ev.(type) {
	case *views.EventWidgetContent:
		p.layout()
		p.PostEventWidgetContent(p)
		return true
	}

	if p.selected != nil {
		if p.selected.HandleEvent(ev) {
			return true
		}
	}

	if ev, ok := ev.(*tcell.EventKey); ok {
		if p.selected != nil {
			switch ev.Key() {
			case tcell.KeyEsc:
				p.unselect()
				return true
			case tcell.KeyTab:
				p.selectNext()
				return true
			}
		} else {
			switch ev.Key() {
			case tcell.KeyTab:
				p.selectNext()
				return true
			}
		}
		return false
	}

	for _, c := range p.children {
		if c.widget != p.selected && c.widget.HandleEvent(ev) {
			return true
		}
	}

	return false
}

func (p *panels) afterAdd(w SelectableWidget) {
	if p.selected == nil && len(p.children) == 1 {
		p.selectIndex(0)
	}
}

func (p *panels) selectNext() {

	if p.selected == nil {
		if len(p.children) > 0 {
			p.selectIndex(0)
		}
		return
	}

	p.selected.SetTheme(theme.ThemeInactive)

	for idx, child := range p.children {
		if p.selected == child.widget {
			idx = (idx + 1) % len(p.children)
			p.selectIndex(idx)
			return
		}
	}
}

func (p *panels) unselect() {
	if p.selected == nil {
		return
	}
	p.selected.SetTheme(theme.ThemeInactive)
	p.selected = nil
}

func (p *panels) selectIndex(idx int) {
	p.selected = p.children[idx].widget.(SelectableWidget)
	p.selected.SetTheme(theme.ThemeActive)
}
