package table

type TR interface {
	ID() string
	Columns() []TD
}

type tableRow struct {
	id   string
	cols []TD
}

func NewTR(id string, cols []TD) TR {
	return &tableRow{id, cols}
}

func (row *tableRow) ID() string {
	return row.id
}

func (row *tableRow) Columns() []TD {
	return row.cols
}
