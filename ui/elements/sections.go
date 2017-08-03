package elements

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/util"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type SelectableWidget interface {
	views.Widget
	theme.Themeable
}

type Sections interface {
	views.Widget

	Widgets() []SelectableWidget

	Append(SelectableWidget)
	Prepend(SelectableWidget)
	Remove(SelectableWidget)
	InsertBefore(SelectableWidget, SelectableWidget)
	InsertAfter(SelectableWidget, SelectableWidget)

	Selected() SelectableWidget
}

type sections struct {
	panes
	selected SelectableWidget
}

func NewVSections(env util.Env, expand bool) Sections {
	return NewSections(env, views.Vertical, expand)
}

func NewHSections(env util.Env, expand bool) Sections {
	return NewSections(env, views.Horizontal, expand)
}

func NewSections(env util.Env, o views.Orientation, expand bool) Sections {
	return &sections{
		panes: panes{orientation: o, expand: expand, env: env},
	}
}

func (p *sections) Selected() SelectableWidget {
	return p.selected
}

func (p *sections) Widgets() []SelectableWidget {
	children := make([]SelectableWidget, 0, len(p.children))
	for _, c := range p.children {
		children = append(children, c.widget.(SelectableWidget))
	}
	return children
}

func (p *sections) Append(w SelectableWidget) {
	p.panes.Append(w)
	p.afterAdd(w)
}

func (p *sections) Prepend(w SelectableWidget) {
	p.panes.Prepend(w)
	p.afterAdd(w)
}

func (p *sections) Remove(w SelectableWidget) {
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

func (p *sections) InsertBefore(mark SelectableWidget, w SelectableWidget) {
	p.panes.InsertBefore(mark, w)
}

func (p *sections) InsertAfter(mark SelectableWidget, w SelectableWidget) {
	p.panes.InsertAfter(mark, w)
}

func (p *sections) HandleEvent(ev tcell.Event) bool {

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

func (p *sections) afterAdd(w SelectableWidget) {
	if p.selected == nil && len(p.children) == 1 {
		p.selectIndex(0)
	}
	if p.selected != w {
		w.SetTheme(theme.ThemeInactive)
	}
}

func (p *sections) selectNext() {

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

func (p *sections) unselect() {
	if p.selected == nil {
		return
	}
	p.selected.SetTheme(theme.ThemeInactive)
	p.selected = nil
}

func (p *sections) selectIndex(idx int) {
	p.selected = p.children[idx].widget.(SelectableWidget)
	p.selected.SetTheme(theme.ThemeActive)
}
