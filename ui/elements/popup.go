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

type Popup struct {
	view views.View

	content views.Widget
	cview   *views.ViewPort

	style tcell.Style

	closer PopupCloser

	xoff int
	yoff int

	width  int
	height int

	ctx Context

	views.WidgetWatchers
}

func NewPopup(ctx Context, width int, height int, style tcell.Style) *Popup {

	w := &Popup{
		closer: KeyEscPopupCloser(),
		width:  width,
		height: height,
		style:  style,
		cview:  views.NewViewPort(nil, 0, 0, 0, 0),
		ctx:    ctx.NewWithID("ui/elements/popup"),
	}
	return w
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

func (p *Popup) SetCloser(closer PopupCloser) {
	p.closer = closer
}

func (p *Popup) HandleEvent(ev tcell.Event) bool {

	switch ev := ev.(type) {
	case *EventPopupClose:
		if ev.Widget() == nil {
			p.Close()
			return true
		}
	}

	if p.closer != nil && p.closer.HandleEvent(ev) {
		p.Close()
		return true
	}

	if p.content != nil {
		return p.content.HandleEvent(ev)
	}
	return false
}

func (p *Popup) Size() (int, int) {
	return p.width, p.height
}

func (p *Popup) Draw() {

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

func (p *Popup) Close() {
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
