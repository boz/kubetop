package table

import (
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/util"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Display interface {
	ResetRows(rows []TR)
	InsertRow(row TR)
	UpdateRow(row TR)
	RemoveRow(id string)
}

type Widget struct {
	model model
	colsz []int

	width  int
	height int
	expand bool

	view  views.View
	hport *views.ViewPort
	rport *views.ViewPort
	views.WidgetWatchers

	env util.Env
}

func NewWidget(env util.Env, cols []TH, expand bool) *Widget {
	return &Widget{
		model:  newModel(cols),
		hport:  views.NewViewPort(nil, 0, 0, 0, 0),
		rport:  views.NewViewPort(nil, 0, 1, 0, 0),
		expand: expand,
		env:    env,
	}
}

func (tw *Widget) ResetRows(rows []TR) {
	tw.model.reset(rows)
	tw.PostEventWidgetContent(tw)
}

func (tw *Widget) InsertRow(row TR) {
	tw.model.insert(row)
	tw.PostEventWidgetContent(tw)
}

func (tw *Widget) UpdateRow(row TR) {
	tw.model.update(row)
	tw.PostEventWidgetContent(tw)
}

func (tw *Widget) RemoveRow(id string) {
	tw.model.remove(id)
	tw.PostEventWidgetContent(tw)
}

func (tw *Widget) Draw() {
	tw.hport.Fill(' ', theme.Base)
	tw.rport.Fill(' ', theme.Base)
	tw.drawHeader()
	tw.model.each(func(roff int, row TR) {
		tw.drawRow(roff, row)
	})
}

func (tw *Widget) Resize() {
	if tw.view == nil {
		return
	}
	tw.resizeContent()
}

func (tw *Widget) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {

		case tcell.KeyUp, tcell.KeyCtrlP:
			return tw.keyUp()
		case tcell.KeyDown, tcell.KeyCtrlN:
			return tw.keyDown()
		case tcell.KeyLeft, tcell.KeyCtrlB:
			return tw.keyLeft()
		case tcell.KeyRight, tcell.KeyCtrlF:
			return tw.keyRight()
		case tcell.KeyEscape:
			return tw.keyEscape()

		case tcell.KeyRune:
			switch ev.Rune() {
			case 'k':
				return tw.keyUp()
			case 'j':
				return tw.keyDown()
			case 'h':
				return tw.keyLeft()
			case 'l':
				return tw.keyRight()
			}
		}
	}
	return false
}

func (tw *Widget) SetView(view views.View) {
	tw.view = view
	tw.hport.SetView(view)
	tw.rport.SetView(view)
	tw.Resize()
}

func (tw *Widget) Size() (int, int) {
	return tw.width, tw.height
}

func (tw *Widget) resizeContent() {

	colsz := make([]int, len(tw.model.columns()))
	update := func(i int, col TD) {
		width, _ := col.Size()
		if width+styleColPad > colsz[i] {
			colsz[i] = width + styleColPad
		}
	}

	for i, col := range tw.model.columns() {
		update(i, col)
	}
	tw.model.each(func(_ int, row TR) {
		for i, col := range row.Columns() {
			update(i, col)
		}
	})

	width := 0
	for _, col := range colsz {
		width += col
	}

	vwidth, _ := tw.view.Size()

	if tw.expand && vwidth > width {
		delta := vwidth - width
		pad := delta / len(colsz)
		rem := delta % len(colsz)

		for idx := range colsz {
			colsz[idx] += pad
			if len(colsz)-idx-1 < rem {
				colsz[idx] += 1
			}
		}
		width = vwidth
	}

	mheight := tw.model.size()

	tw.hport.Resize(0, 0, width, 1)
	tw.hport.SetContentSize(width, 1, false)

	tw.rport.Resize(0, 1, width, mheight)
	tw.rport.SetContentSize(width, mheight, false)

	tw.colsz = colsz
	tw.width = width
	tw.height = mheight + 1

	tw.scrollToActive()
}

func (tw *Widget) drawHeader() {
	xoff := 0
	yoff := 0
	cols := tw.model.columns()
	view := tw.hport
	for i, col := range cols {
		width := tw.colsz[i]
		cview := newView(view, xoff, yoff, width, 1, theme.Table.TH)
		col.Draw(cview)
		xoff += width
	}
}

func (tw *Widget) drawRow(yoff int, row TR) {
	xoff := 0
	cols := row.Columns()
	view := tw.rport

	lth := theme.Table.TD
	if tw.model.isActive(row.ID()) {
		lth = theme.Table.TDSelected
	}

	for i, col := range cols {
		width := tw.colsz[i]
		cview := newView(view, xoff, yoff, width, 1, lth)
		col.Draw(cview)
		xoff += width
	}
}

func (tw *Widget) keyUp() bool {
	if tw.model.activatePrev() {
		tw.scrollToActive()
		tw.postActive()
		return true
	}
	return false
}

func (tw *Widget) keyDown() bool {
	if tw.model.activateNext() {
		tw.scrollToActive()
		tw.postActive()
		return true
	}
	return false
}

func (tw *Widget) keyLeft() bool {
	tw.hport.ScrollLeft(1)
	tw.rport.ScrollLeft(1)
	return true
}
func (tw *Widget) keyRight() bool {
	tw.hport.ScrollRight(1)
	tw.rport.ScrollRight(1)
	return true
}

func (tw *Widget) keyEscape() bool {
	if tw.model.clearActive() {
		tw.PostEvent(newEventRowInactive(tw))
		return true
	}
	return false
}

func (tw *Widget) scrollToActive() {
	if idx, _ := tw.model.getActive(); idx >= 0 {
		tw.rport.MakeVisible(-1, idx)
	}
}

func (tw *Widget) postActive() {
	if idx, row := tw.model.getActive(); idx >= 0 {
		tw.PostEvent(newEventRowActive(tw, row))
	}
}

func (tw *Widget) _debug(ctx string) {
	w, h := tw.view.Size()
	tw.env.Log().Debugf("%v: view.Size():  (%v,%v)", ctx, w, h)

	w, h = tw.hport.Size()
	tw.env.Log().Debugf("%v: hport.Size(): (%v,%v)", ctx, w, h)

	w, h = tw.rport.Size()
	tw.env.Log().Debugf("%v: rport.Size(): (%v,%v)", ctx, w, h)

	tw.env.Log().Debugf("%v: model.size():  %v", ctx, tw.model.size())
}
