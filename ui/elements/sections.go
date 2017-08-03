package elements

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/util"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Styleable interface {
	views.Widget
	SetStyle(tcell.Style)
}

type Sections interface {
	views.Widget

	Widgets() []Themeable

	Append(Themeable)
	Prepend(Themeable)
	Remove(Themeable)
	InsertBefore(Themeable, Themeable)
	InsertAfter(Themeable, Themeable)

	Selected() Themeable
}

type sections struct {
	panes
	selected Themeable
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

func (p *sections) Selected() Themeable {
	return p.selected
}

func (p *sections) Widgets() []Themeable {
	children := make([]Themeable, 0, len(p.children))
	for _, c := range p.children {
		children = append(children, c.widget.(Themeable))
	}
	return children
}

func (p *sections) Append(w Themeable) {
	p.panes.Append(w)
	p.afterAdd(w)
}

func (p *sections) Prepend(w Themeable) {
	p.panes.Prepend(w)
	p.afterAdd(w)
}

func (p *sections) Remove(w Themeable) {
	if p.selected == nil || p.selected != w {
		p.panes.Remove(w)
		return
	}
	for idx, c := range p.children {
		if p.selected != c.widget.(Themeable) {
			continue
		}
		if idx >= len(p.children)-1 {
			continue
		}

		wnext := p.children[idx+1].widget.(Themeable)

		if wnext == p.selected {
			continue
		}

		p.selected = wnext
		p.selected.SetTheme(theme.ThemeActive)
		break
	}
	p.panes.Remove(w)
}

func (p *sections) InsertBefore(mark Themeable, w Themeable) {
	p.panes.InsertBefore(mark, w)
}

func (p *sections) InsertAfter(mark Themeable, w Themeable) {
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

func (p *sections) afterAdd(w Themeable) {
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
	p.selected = p.children[idx].widget.(Themeable)
	p.selected.SetTheme(theme.ThemeActive)
}
