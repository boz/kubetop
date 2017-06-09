package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	styleTableTH = tcell.StyleDefault.Bold(true)

	tableColPad = 2
)

type CellView interface {
	Size() (int, int)
	SetText(string, tcell.Style)
	SetContent(x, y int, ch rune, comb []rune, s tcell.Style)
}

type cellView struct {
	view     views.View
	xoff     int
	yoff     int
	width    int
	height   int
	selected bool
}

func NewCellView(view views.View, xoff, yoff, width, height int, selected bool) CellView {
	return &cellView{view, xoff, yoff, width, height, selected}
}

func (v *cellView) Size() (int, int) {
	return v.width, v.height
}

func (v *cellView) SetContent(x, y int, ch rune, comb []rune, s tcell.Style) {
	v.view.SetContent(x+v.xoff, y+v.yoff, ch, comb, s.Reverse(v.selected))
}

func (v *cellView) SetText(text string, s tcell.Style) {
	for i, ch := range text {
		v.SetContent(i, 0, ch, nil, s)
	}
}

type Table interface {
	Header() TableHeader
	Rows() []TableRow
	AddRow(TableRow)
	RemoveRow(string)
}

type TableHeader interface {
	Columns() []TableColumn
}

type TableRow interface {
	ID() string
	Columns() []TableColumn
}

type TableColumn interface {
	ID() string
	Size() (int, int)
	Draw(CellView)
}

type TableWidget struct {
	model Table
	view  views.View
	colsz []int
	port  *views.ViewPort
	views.WidgetWatchers

	curRow string
}

func NewTableWidget(model Table) *TableWidget {
	return &TableWidget{
		model: model,
		port:  views.NewViewPort(nil, 0, 0, 0, 0),
	}
}

func (tw *TableWidget) Draw() {
	tw.view.Fill(' ', tcell.StyleDefault)
	tw.drawHeader()
	for roff, row := range tw.model.Rows() {
		tw.drawRow(roff, row)
	}
}

func (tw *TableWidget) Resize() {
	if tw.view == nil {
		return
	}
	tw.resizeContent()
}

func (tw *TableWidget) resizeContent() {
	colsz := make([]int, len(tw.model.Header().Columns()))
	tw.adjustColSizesForRow(colsz, tw.model.Header().Columns())
	for _, row := range tw.model.Rows() {
		tw.adjustColSizesForRow(colsz, row.Columns())
	}

	width := 0
	for _, col := range colsz {
		width += col
	}
	height := len(tw.model.Rows())

	tw.port.Resize(0, 1, width, height)

	tw.colsz = colsz
}

func (tw *TableWidget) HandleEvent(ev tcell.Event) bool {
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

func (tw *TableWidget) SetView(view views.View) {
	tw.view = view
	tw.port.SetView(view)
	tw.Resize()
}

func (tw *TableWidget) Size() (int, int) {
	x, y := tw.port.Size()
	return x, y + 1
}

func (tw *TableWidget) adjustColSizesForRow(widths []int, cols []TableColumn) {
	for i, col := range cols {
		w, _ := col.Size()
		if w+tableColPad > widths[i] {
			widths[i] = w + tableColPad
		}
	}
}

func (tw *TableWidget) drawHeader() {
	xoff := 0
	yoff := 0
	cols := tw.model.Header().Columns()
	view := tw.view
	for i, col := range cols {
		width := tw.colsz[i]
		cview := NewCellView(view, xoff, yoff, width, 1, false)
		col.Draw(cview)
		xoff += width
	}
}

func (tw *TableWidget) drawRow(yoff int, row TableRow) {
	xoff := 0
	cols := row.Columns()
	view := tw.port
	selected := tw.curRow == row.ID()
	for i, col := range cols {
		width := tw.colsz[i]
		cview := NewCellView(view, xoff, yoff, width, 1, selected)
		col.Draw(cview)
		xoff += width
	}
}

func (tw *TableWidget) keyUp() bool {
	curidx, ok := tw.currentRowIndex()
	if !ok {
		return false
	}
	curidx -= 1
	if curidx < 0 {
		return true
	}
	rows := tw.model.Rows()
	tw.curRow = rows[curidx].ID()
	tw.port.MakeVisible(-1, curidx)
	return true
}

func (tw *TableWidget) keyDown() bool {
	curidx, ok := tw.currentRowIndex()
	if !ok {
		curidx = -1
	}
	curidx += 1
	rows := tw.model.Rows()

	if curidx >= len(rows) {
		return true
	}

	tw.curRow = rows[curidx].ID()
	tw.port.MakeVisible(-1, curidx)
	return true
}
func (tw *TableWidget) keyLeft() bool {
	tw.port.ScrollLeft(1)
	return true
}
func (tw *TableWidget) keyRight() bool {
	tw.port.ScrollRight(1)
	return true
}
func (tw *TableWidget) keyEscape() bool {
	return true
}

func (tw *TableWidget) currentRowIndex() (int, bool) {
	if tw.curRow == "" {
		return 0, false
	}
	for i, row := range tw.model.Rows() {
		if row.ID() == tw.curRow {
			return i, true
		}
	}
	return 0, false
}

type tableModel struct {
	header TableHeader
	rows   []TableRow
}

func (t *tableModel) Header() TableHeader {
	return t.header
}

func (t *tableModel) Rows() []TableRow {
	return t.rows
}

func (t *tableModel) AddRow(new TableRow) {
	for i, row := range t.rows {
		if row.ID() == new.ID() {
			t.rows[i] = new
			return
		}
	}
	t.rows = append(t.rows, new)
}

func (t *tableModel) RemoveRow(id string) {
	rows := make([]TableRow, 0, len(t.rows))
	for _, row := range t.rows {
		if row.ID() != id {
			rows = append(rows, row)
		}
	}
	t.rows = rows
}

func NewTable(header TableHeader, rows []TableRow) Table {
	return &tableModel{header, rows}
}

type tableHeader struct {
	cols []TableColumn
}

func NewTableHeader(cols []TableColumn) TableHeader {
	return &tableHeader{cols}
}

func (th *tableHeader) Columns() []TableColumn {
	return th.cols
}

type tableColumn struct {
	id    string
	text  string
	style tcell.Style
}

func NewTableColumn(id string, text string, style tcell.Style) TableColumn {
	return &tableColumn{id, text, style}
}

func (col *tableColumn) ID() string {
	return col.id
}

func (col *tableColumn) Size() (int, int) {
	return len(col.text), 1
}

func (col *tableColumn) Draw(view CellView) {
	view.SetText(col.text, col.style)
}

type tableRow struct {
	id   string
	cols []TableColumn
}

func NewTableTH(id string, text string) TableColumn {
	return NewTableColumn(id, text, styleTableTH)
}

func NewTableRow(id string, cols []TableColumn) TableRow {
	return &tableRow{id, cols}
}

func (row *tableRow) ID() string {
	return row.id
}

func (row *tableRow) Columns() []TableColumn {
	return row.cols
}
