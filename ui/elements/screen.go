package elements

import (
	"github.com/gdamore/tcell/views"
)

type Screen interface {
	Widget

	State() ScreenState

	Help() []views.Widget
}

type screen struct {
	widget
	state ScreenState
	help  []views.Widget
}

func NewScreen(ctx Context, req Request, title string, content views.Widget, help []views.Widget) Screen {
	return &screen{widget{content: content, ctx: ctx}, NewScreenState(req, title), help}
}

func (s *screen) State() ScreenState {
	return s.state
}

func (s *screen) Help() []views.Widget {
	return s.help
}

type ScreenState interface {
	Request() Request
	Title() string
}

func NewScreenState(req Request, title string) ScreenState {
	return screenState{req, title}
}

type screenState struct {
	request Request
	title   string
}

func (ss screenState) Request() Request {
	return ss.request
}

func (ss screenState) Title() string {
	return ss.title
}
