package table

import (
	"container/list"
	"sort"
	"strings"
)

type model interface {
	columns() []TH
	each(func(int, TR))
	size() int

	reset([]TR)
	insert(TR)
	update(TR)
	remove(string)

	activateNext() bool
	activatePrev() bool
	isActive(string) bool
	getActive() (int, TR)
	clearActive() bool
}

type _model struct {
	cols     []TH
	sortcols []int
	rows     *list.List
	active   *list.Element
}

func newModel(cols []TH) *_model {
	model := &_model{
		cols: cols,
		rows: list.New(),
	}
	model.sortPrep()
	return model
}

func (m *_model) columns() []TH {
	return m.cols
}

func (m *_model) size() int {
	return m.rows.Len()
}

func (m *_model) each(fn func(int, TR)) {
	i := 0
	for e := m.rows.Front(); e != nil; e = e.Next() {
		fn(i, e.Value.(TR))
		i++
	}
}

func (m *_model) reset(rows []TR) {
	sort.Slice(rows, func(i, j int) bool {
		return m.compare(rows[i], rows[j]) < 0
	})
	m.rows.Init()
	for _, row := range rows {
		m.rows.PushBack(row)
	}
}

func (m *_model) insert(row TR) {
	m.doInsert(row)
}

func (m *_model) update(row TR) {
	cur := m.find(row.ID())
	if cur == nil {
		m.doInsert(row)
		return
	}

	if m.compare(row, cur.Value.(TR)) == 0 {
		cur.Value = row
		return
	}

	m.rows.Remove(cur)
	next := m.doInsert(row)

	if m.active == cur {
		m.active = next
	}
}

func (m *_model) remove(id string) {
	if e := m.find(id); e != nil {
		if e == m.active {
			if m.active = e.Next(); m.active == nil {
				m.active = e.Prev()
			}
		}
		m.rows.Remove(e)
	}
}

func (m *_model) activateNext() bool {
	if m.active == nil {
		m.active = m.rows.Front()
		return m.active != nil
	}
	m.active = m.active.Next()
	if m.active == nil {
		m.active = m.rows.Front()
	}
	return true
}

func (m *_model) activatePrev() bool {
	if m.active == nil {
		return false
	}
	m.active = m.active.Prev()
	if m.active == nil {
		m.active = m.rows.Back()
	}
	return true
}

func (m *_model) isActive(id string) bool {
	if m.active == nil {
		return false
	}
	return m.active.Value.(TR).ID() == id
}

func (m *_model) getActive() (int, TR) {
	if m.active == nil {
		return -1, nil
	}
	return m.elIndex(m.active), m.active.Value.(TR)
}

func (m *_model) clearActive() bool {
	if m.active == nil {
		return false
	}
	m.active = nil
	return true
}

func (m *_model) doInsert(row TR) *list.Element {
	for e := m.rows.Front(); e != nil; e = e.Next() {
		if m.compare(row, e.Value.(TR)) < 0 {
			return m.rows.InsertBefore(row, e)
		}
	}
	return m.rows.PushBack(row)
}

func (m *_model) find(id string) *list.Element {
	for e := m.rows.Front(); e != nil; e = e.Next() {
		if id == e.Value.(TR).ID() {
			return e
		}
	}
	return nil
}

func (m *_model) sortPrep() {
	sortcols := make([]int, 0)
	for idx, col := range m.cols {
		if sidx := col.SortOrder(); sidx >= 0 && col.Sortable() {
			sortcols = append(sortcols, idx)
		}
	}
	sort.SliceStable(sortcols, func(i, j int) bool {
		return m.cols[i].SortOrder() < m.cols[j].SortOrder()
	})
	m.sortcols = sortcols
}

func (m *_model) compare(a, b TR) int {
	acols, bcols := a.Columns(), b.Columns()
	for _, idx := range m.sortcols {
		if val := strings.Compare(acols[idx].Key(), bcols[idx].Key()); val != 0 {
			return val
		}
	}
	return 0
}

func (m *_model) elIndex(el *list.Element) int {
	var e *list.Element
	var idx int
	for e = m.rows.Front(); el != nil && e != el && e != nil; e = e.Next() {
		idx++
	}
	if e == nil {
		return -1
	}
	return idx
}
