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
	view   views.View
	xoff   int
	yoff   int
	width  int
	height int
}

func NewCellView(view views.View, xoff, yoff, width, height int) CellView {
	return &cellView{view, xoff, yoff, width, height}
}

func (v *cellView) Size() (int, int) {
	return v.width, v.height
}

func (v *cellView) SetContent(x, y int, ch rune, comb []rune, s tcell.Style) {
	v.view.SetContent(x+v.xoff, y+v.yoff, ch, comb, s)
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
}

func NewTableWidget(model Table) *TableWidget {
	return &TableWidget{
		model: model,
		port:  views.NewViewPort(nil, 0, 0, 0, 0),
	}
}

func (tw *TableWidget) Draw() {
	tw.view.Fill(' ', tcell.StyleDefault)
	tw.drawRow(tw.view, 0, tw.model.Header().Columns())
	for roff, row := range tw.model.Rows() {
		tw.drawRow(tw.port, roff, row.Columns())
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
	length := len(tw.model.Rows())

	tw.port.SetContentSize(width, length, false)
	tw.port.Resize(0, 1, width, length)

	tw.colsz = colsz
}

func (tw *TableWidget) HandleEvent(ev tcell.Event) bool {
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

func (tw *TableWidget) drawRow(view views.View, yoff int, cols []TableColumn) {
	xoff := 0
	for i, col := range cols {
		width := tw.colsz[i]
		cview := NewCellView(view, xoff, yoff, width, 1)
		col.Draw(cview)
		xoff += width
	}
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
