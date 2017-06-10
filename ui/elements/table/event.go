package table

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type RowEvent interface {
	views.EventWidget
	Row() TR
}

type EventRowActive struct {
	rowEvent
}

type EventRowInactive struct {
	rowEvent
}

type EventRowSelected struct {
	rowEvent
}

func newEventRowActive(widget *Widget, row TR) RowEvent {
	ev := &EventRowActive{rowEvent: rowEvent{widget: widget, row: row}}
	ev.SetEventNow()
	return ev
}

func newEventRowInactive(widget *Widget) RowEvent {
	ev := &EventRowInactive{rowEvent: rowEvent{widget: widget}}
	ev.SetEventNow()
	return ev
}

func newEventRowSelected(widget *Widget, row TR) RowEvent {
	ev := &EventRowSelected{rowEvent: rowEvent{widget: widget, row: row}}
	ev.SetEventNow()
	return ev
}

type rowEvent struct {
	widget *Widget
	row    TR
	tcell.EventTime
}

func (ev *rowEvent) Widget() views.Widget {
	return ev.widget
}
func (ev *rowEvent) Row() TR {
	return ev.row
}
