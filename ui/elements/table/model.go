package table

import (
	"container/list"
	"sort"
	"strings"
)

type model interface {
	columns() []TH
	each(func(int, TR))
	reset([]TR)
	insert(TR)
	update(TR)
	remove(string)
	selectNext() (int, TR)
	selectPrev() (int, TR)
	isSelected(string) bool
	clearSelection() bool
}

type _model struct {
	cols     []TH
	sortcols []int
	rows     *list.List
	selected *list.Element
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

func (m *_model) each(fn func(int, TR)) {
	i := 0
	for e := m.rows.Front(); e != nil; e = e.Next() {
		row := e.Value.(TR)
		fn(i, row)
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
	for e := m.rows.Front(); e != nil; e = e.Next() {
		if m.compare(row, e.Value.(TR)) < 0 {
			m.rows.InsertBefore(row, e)
			return
		}
	}
	m.rows.PushBack(row)
}

func (m *_model) update(row TR) {
	cur := m.find(row.ID())

	if cur == nil {
		m.insert(row)
		return
	}

	cur.Value = row

	for e := cur.Next(); e != nil; e = e.Next() {
		if m.compare(row, e.Value.(TR)) >= 0 {
			break
		}
		m.rows.MoveBefore(cur, e)
	}

	for e := cur.Prev(); e != nil; e = e.Prev() {
		if m.compare(row, e.Value.(TR)) <= 0 {
			break
		}
		m.rows.MoveAfter(cur, e)
	}
}

func (m *_model) remove(id string) {
	if e := m.find(id); e != nil {
		if e == m.selected {
			m.selected = m.selected.Prev()
		}
		m.rows.Remove(e)
	}
}

func (m *_model) selectNext() (int, TR) {
	if m.selected == nil {
		m.selected = m.rows.Front()
		goto done
	}
	if m.selected.Next() == nil {
		m.selected = m.rows.Front()
		goto done
	}
	m.selected = m.selected.Next()

done:
	if m.selected == nil {
		return -1, nil
	}
	return m.elIndex(m.selected), m.selected.Value.(TR)
}

func (m *_model) selectPrev() (int, TR) {
	if m.selected == nil {
		return -1, nil
	}
	if m.selected.Prev() == nil {
		m.selected = m.rows.Back()
		goto done
	}
	m.selected = m.selected.Prev()

done:
	if m.selected == nil {
		return -1, nil
	}
	return m.elIndex(m.selected), m.selected.Value.(TR)
}

func (m *_model) isSelected(id string) bool {
	if m.selected == nil {
		return false
	}
	return m.selected.Value.(TR).ID() == id
}

func (m *_model) clearSelection() bool {
	if m.selected == nil {
		return false
	}
	m.selected = nil
	return true
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
