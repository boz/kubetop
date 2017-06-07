package elements

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
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

	widths := make([]int, len(tw.model.Header().Columns()))

	tw.adjustColSizesForRow(widths, tw.model.Header().Columns())
	for _, row := range tw.model.Rows() {
		tw.adjustColSizesForRow(widths, row.Columns())
	}

	tw.drawRow(widths, 0, tw.model.Header().Columns())
	for roff, row := range tw.model.Rows() {
		tw.drawRow(widths, roff+1, row.Columns())
	}
}

func (tw *TableWidget) Resize() {
}

func (tw *TableWidget) HandleEvent(ev tcell.Event) bool {
	return false
}

func (tw *TableWidget) SetView(view views.View) {
	tw.view = view
	tw.port.SetView(view)
}

func (tw *TableWidget) Size() (int, int) {
	return tw.view.Size()
}

func (tw *TableWidget) adjustColSizesForRow(widths []int, cols []TableColumn) {
	for i, col := range cols {
		w, _ := col.Size()
		if w > widths[i] {
			widths[i] = w
		}
	}
}

func (tw *TableWidget) drawRow(widths []int, yoff int, cols []TableColumn) {
	xoff := 0
	for i, col := range cols {
		width := widths[i]
		view := NewCellView(tw.view, xoff, yoff, width, 1)
		col.Draw(view)
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
	pad   int
}

func NewTableColumn(id string, text string, style tcell.Style, pad int) TableColumn {
	return &tableColumn{id, text, style, pad}
}

func (col *tableColumn) ID() string {
	return col.id
}

func (col *tableColumn) Size() (int, int) {
	return len(col.text) + col.pad, 1
}

func (col *tableColumn) Draw(view CellView) {
	view.SetText(col.text, col.style)
}

type tableRow struct {
	id   string
	cols []TableColumn
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
