package elements

import "github.com/gdamore/tcell/views"

type Screen interface {
	Widget

	State() ScreenState
}

type screen struct {
	widget
	state ScreenState
}

func NewScreen(ctx Context, req Request, title string, content views.Widget) Screen {
	return &screen{widget{content, ctx}, NewScreenState(req, title)}
}

func (s *screen) State() ScreenState {
	return s.state
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
