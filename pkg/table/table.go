package table

import (
	"sync"
)

type Table[R comparable, C comparable, V any] struct {
	data         map[R]map[C]V
	sync.RWMutex // 如果需要线程安全
}

func NewTable[R comparable, C comparable, V any]() *Table[R, C, V] {
	return &Table[R, C, V]{
		data: make(map[R]map[C]V),
	}
}

func (t *Table[R, C, V]) Put(row R, col C, value V) {
	t.Lock()
	defer t.Unlock()

	if t.data[row] == nil {
		t.data[row] = make(map[C]V)
	}
	t.data[row][col] = value
}

func (t *Table[R, C, V]) Get(row R, col C) (V, bool) {
	t.RLock()
	defer t.RUnlock()

	colMap, rowExists := t.data[row]
	if !rowExists {
		var zero V
		return zero, false
	}
	value, colExists := colMap[col]
	return value, colExists
}

func (t *Table[R, C, V]) Row(row R) map[C]V {
	t.RLock()
	defer t.RUnlock()

	// 返回副本以避免并发问题
	rowData := make(map[C]V)
	for k, v := range t.data[row] {
		rowData[k] = v
	}
	return rowData
}

func (t *Table[R, C, V]) Rows() []R {
	t.RLock()
	defer t.RUnlock()

	keys := make([]R, 0, len(t.data))
	for r := range t.data {
		keys = append(keys, r)
	}
	return keys
}

// Size 新增：返回当前表的行数
func (t *Table[R, C, V]) Size() int {
	t.RLock()
	defer t.RUnlock()
	return len(t.data)
}

func (t *Table[R, C, V]) Column(col C) map[R]V {
	t.RLock()
	defer t.RUnlock()

	columnData := make(map[R]V)
	for row, colMap := range t.data {
		if value, exists := colMap[col]; exists {
			columnData[row] = value
		}
	}
	return columnData
}

func (t *Table[R, C, V]) Remove(row R, col C) bool {
	t.Lock()
	defer t.Unlock()

	if colMap, exists := t.data[row]; exists {
		if _, colExists := colMap[col]; colExists {
			delete(colMap, col)
			// 如果该行已空，删除整行
			if len(colMap) == 0 {
				delete(t.data, row)
			}
			return true
		}
	}
	return false
}

func (t *Table[R, C, V]) Clear() {
	t.Lock()
	defer t.Unlock()

	t.data = make(map[R]map[C]V)
}
