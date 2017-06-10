package table

type TH interface {
	TD
	Sortable() bool
	SortOrder() int
}

type tableTH struct {
	tableTD
	sortable  bool
	sortOrder int
}

func NewTH(id string, text string, sortable bool, sortOrder int) TH {
	return &tableTH{
		tableTD: tableTD{
			id:    id,
			text:  text,
			style: styleTH,
		},
		sortable:  sortable,
		sortOrder: sortOrder,
	}
	return nil
}

func (th *tableTH) Sortable() bool {
	return th.sortable
}

func (th *tableTH) SortOrder() int {
	return th.sortOrder
}
