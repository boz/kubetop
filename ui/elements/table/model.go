package table

import (
	"container/list"
	"sort"
	"strings"
)

type tableModel struct {
	cols     []TH
	sortcols []int
	rows     *list.List
	selected *list.Element
}

func newTableModel(cols []TH) *tableModel {
	model := &tableModel{
		cols: cols,
		rows: list.New(),
	}
	model.sortPrep()
	return model
}

func (m *tableModel) columns() []TH {
	return m.cols
}

func (m *tableModel) each(fn func(int, TR)) {
	i := 0
	for e := m.rows.Front(); e != nil; e = e.Next() {
		row := e.Value.(TR)
		fn(i, row)
		i += 1
	}
}

func (m *tableModel) reset(rows []TR) {
	sort.Slice(rows, func(i, j int) bool {
		return m.compare(rows[i], rows[j]) < 0
	})
	m.rows.Init()
	for _, row := range rows {
		m.rows.PushBack(row)
	}
}

func (m *tableModel) insert(row TR) {
	for e := m.rows.Front(); e != nil; e = e.Next() {
		if m.compare(row, e.Value.(TR)) < 0 {
			m.rows.InsertBefore(row, e)
			return
		}
	}
	m.rows.PushBack(row)
}

func (m *tableModel) update(row TR) {
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

func (m *tableModel) remove(id string) {
	if e := m.find(id); e != nil {
		if e == m.selected {
			m.selected = m.selected.Prev()
		}
		m.rows.Remove(e)
	}
}

func (m *tableModel) selectNext() (int, TR) {
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

func (m *tableModel) selectPrev() (int, TR) {
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

func (m *tableModel) isSelected(id string) bool {
	if m.selected == nil {
		return false
	}
	return m.selected.Value.(TR).ID() == id
}

func (m *tableModel) clearSelection() bool {
	if m.selected == nil {
		return false
	}
	m.selected = nil
	return true
}

func (m *tableModel) find(id string) *list.Element {
	for e := m.rows.Front(); e != nil; e = e.Next() {
		if id == e.Value.(TR).ID() {
			return e
		}
	}
	return nil
}

func (m *tableModel) sortPrep() {
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

func (m *tableModel) compare(a, b TR) int {
	acols, bcols := a.Columns(), b.Columns()
	for _, idx := range m.sortcols {
		if val := strings.Compare(acols[idx].Key(), bcols[idx].Key()); val != 0 {
			return val
		}
	}
	return 0
}

func (m *tableModel) elIndex(el *list.Element) int {
	var e *list.Element
	var idx int
	for e = m.rows.Front(); el != nil && e != el && e != nil; e = e.Next() {
		idx += 1
	}
	if e == nil {
		return -1
	}
	return idx
}
