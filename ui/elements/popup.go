package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Popup struct {
	view views.View

	content views.Widget
	cview   *views.ViewPort

	style tcell.Style

	xoff int
	yoff int

	width  int
	height int

	views.WidgetWatchers
}

func NewPopup(width int, height int, style tcell.Style) *Popup {
	return &Popup{
		width:  width,
		height: height,
		style:  style,
		cview:  views.NewViewPort(nil, 0, 0, 0, 0),
	}
}

func (p *Popup) SetView(view views.View) {
	p.view = view
	p.cview.SetView(p.view)
	p.Resize()
}

func (p *Popup) SetContent(w views.Widget) {
	p.content = w
	w.SetView(p.cview)
}

func (p *Popup) HandleEvent(ev tcell.Event) bool {
	if p.content != nil {
		return p.content.HandleEvent(ev)
	}
	return false
}

func (p *Popup) Size() (int, int) {
	return p.width, p.height
}

func (p *Popup) Draw() {

	// top
	for x, y := 0, 0; x < p.width; x++ {
		p.view.SetContent(x+p.xoff, y+p.yoff, 'x', nil, p.style)
	}

	// bottom
	for x, y := 0, p.height-1; x < p.width; x++ {
		p.view.SetContent(x+p.xoff, y+p.yoff, 'x', nil, p.style)
	}

	// left
	for x, y := 0, 1; y < p.height-1; y++ {
		p.view.SetContent(x+p.xoff, y+p.yoff, 'x', nil, p.style)
	}

	// right
	for x, y := p.width-1, 1; y < p.height-1; y++ {
		p.view.SetContent(x+p.xoff, y+p.yoff, 'x', nil, p.style)
	}

	if p.content != nil {
		p.content.Draw()
	}
}

func (p *Popup) Resize() {
	if p.view == nil {
		return
	}

	vx, vy := p.view.Size()

	p.xoff, p.yoff = 0, 0

	if xdelta := vx - p.width; xdelta > 1 {
		p.xoff = xdelta / 2
	}

	if ydelta := vy - p.width; ydelta > 1 {
		p.yoff = ydelta / 2
	}

	p.cview.Resize(p.xoff+1, p.yoff+1, p.width-2, p.height-2)

	if p.content != nil {
		p.content.Resize()
	}
}
