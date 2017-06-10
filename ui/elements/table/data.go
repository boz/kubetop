package table

import "github.com/boz/kubetop/ui/theme"

type TD interface {
	ID() string
	Size() (int, int)
	Draw(View)
	Key() string
}

type tableTD struct {
	id   string
	text string
	tv   theme.LabelVariant
}

func NewTD(id string, text string, th theme.LabelVariant) TD {
	return &tableTD{id, text, th}
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

func (col *tableTD) Draw(view View) {
	view.SetText(col.text, col.tv)
}
