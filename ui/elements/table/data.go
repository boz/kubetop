package table

import "github.com/gdamore/tcell"

type TD interface {
	ID() string
	Size() (int, int)
	Draw(CellView)
	Key() string
}

type tableTD struct {
	id    string
	text  string
	style tcell.Style
}

func NewTD(id string, text string, style tcell.Style) TD {
	return &tableTD{id, text, style}
}

func (col *tableTD) ID() string {
	return col.id
}

func (col *tableTD) Size() (int, int) {
	return len(col.text), 1
}

func (col *tableTD) Key() string {
	return col.text
}

func (col *tableTD) Draw(view CellView) {
	view.SetText(col.text, col.style)
}
