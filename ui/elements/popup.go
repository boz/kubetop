package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const (
	boxBorderTopLeft     = '╔'
	boxBorderBottomLeft  = '╚'
	boxBorderTopRight    = '╗'
	boxBorderBottomRight = '╝'
	boxBorderTop         = '═'
	boxBorderRight       = '║'
	boxBorderBottom      = boxBorderTop
	boxBorderLeft        = boxBorderRight
)

type popup struct {
	view views.View

	content views.Widget
	viewp   *views.ViewPort

	style tcell.Style

	closer PopupCloser

	xoff int
	yoff int

	width  int
	height int

	ctx Context

	views.WidgetWatchers
}

func NewPopup(ctx Context, style tcell.Style, content views.Widget) views.Widget {
	w := &popup{
		content: content,
		viewp:   views.NewViewPort(nil, 0, 0, 0, 0),
		style:   style,
		closer:  KeyEscPopupCloser(),
		ctx:     ctx.NewWithID("ui/elements/popup"),
	}
	return w
}

func (p *popup) SetView(view views.View) {
	p.view = view
	p.viewp.SetView(p.view)
	p.content.SetView(p.viewp)
	p.layout()
}

func (p *popup) HandleEvent(ev tcell.Event) bool {

	if p.content.HandleEvent(ev) {
		return true
	}

	if p.closer != nil && p.closer.HandleEvent(ev) {
		p.Close()
		return true
	}

	return false
}

func (p *popup) Size() (int, int) {
	return p.width, p.height
}

func (p *popup) Draw() {

	if p.width <= 0 || p.height <= 0 {
		return
	}

	// top left
	p.view.SetContent(p.xoff, p.yoff, boxBorderTopLeft, nil, p.style)
	// top right
	p.view.SetContent(p.xoff+p.width-1, p.yoff, boxBorderTopRight, nil, p.style)

	// bot right
	p.view.SetContent(p.xoff+p.width-1, p.yoff+p.height-1, boxBorderBottomRight, nil, p.style)
	// bot left
	p.view.SetContent(p.xoff, p.yoff+p.height-1, boxBorderBottomLeft, nil, p.style)

	for x := 1; x < p.width-1; x++ {
		// top
		p.view.SetContent(p.xoff+x, p.yoff, boxBorderTop, nil, p.style)
		// bottom
		p.view.SetContent(p.xoff+x, p.yoff+p.height-1, boxBorderBottom, nil, p.style)
	}

	for y := 1; y < p.height-1; y++ {
		// left
		p.view.SetContent(p.xoff, p.yoff+y, boxBorderLeft, nil, p.style)
		// right
		p.view.SetContent(p.xoff+p.width-1, p.yoff+y, boxBorderRight, nil, p.style)
	}

	for x := 1; x < p.width-1; x++ {
		for y := 1; y < p.height-1; y++ {
			p.view.SetContent(p.xoff+x, p.yoff+y, ' ', nil, p.style)
		}
	}

	if p.content != nil {
		p.content.Draw()
	}
}

func (p *popup) Resize() {
	p.layout()
}

func (p *popup) layout() {
	if p.view == nil {
		return
	}

	vx, vy := p.view.Size()
	px, py := vx, vy
	wx, wy := p.content.Size()

	xoff, yoff := 0, 0

	if delta := vx - wx; delta > 2 {
		xoff = delta / 2
		px = wx + 2
	}

	if delta := vy - wx; delta > 2 {
		yoff = delta / 2
		py = wy + 2
	}

	p.xoff = xoff
	p.yoff = yoff
	p.width = px
	p.height = py

	p.viewp.Resize(p.xoff+1, p.yoff+1, p.width-2, p.height-2)
	p.content.Resize()
}

func (p *popup) Close() {
	ev := &EventPopupClose{}
	ev.SetWidget(p)
	ev.SetEventNow()
	p.PostEvent(ev)
}

type EventPopupClose views.EventWidgetMove

type PopupCloser tcell.EventHandler

func NewPopupCloser(handler func(tcell.Event) bool) PopupCloser {
	return popupCloser{handler}
}

func KeyEscPopupCloser() PopupCloser {
	return NewPopupCloser(func(ev tcell.Event) bool {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEsc:
				return true
			}
		}
		return false
	})
}

type popupCloser struct {
	handler func(ev tcell.Event) bool
}

func (pc popupCloser) HandleEvent(ev tcell.Event) bool {
	return pc.handler(ev)
}
