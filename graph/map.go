package graph

type OrderedMap[K, V comparable] struct {
	keys []K
	data map[K]V
}

func NewOrderedMap[K, V comparable]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys: make([]K, 0),
		data: make(map[K]V),
	}
}

func (m *OrderedMap[K, V]) Set(key K, value V) {
	if _, ok := m.data[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.data[key] = value
}

func (m *OrderedMap[K, V]) Get(key K) V {
	return m.data[key]
}

func (m *OrderedMap[K, V]) Keys() []K {
	return m.keys
}

func (m *OrderedMap[K, V]) Values() []V {
	var values []V
	for _, key := range m.keys {
		values = append(values, m.data[key])
	}
	return values
}

func (m *OrderedMap[K, V]) Len() int {
	return len(m.keys)
}

func (m *OrderedMap[K, V]) Delete(key K) {
	delete(m.data, key)
	for i, k := range m.keys {
		if k == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

func (m *OrderedMap[K, V]) Clear() {
	m.keys = make([]K, 0)
	m.data = make(map[K]V)
}

func (m *OrderedMap[K, V]) Has(key K) bool {
	_, ok := m.data[key]
	return ok
}

func OrderedSetFromSlice[T comparable](s []T) *OrderedSet[T] {
	set := NewOrderedSet[T]()

	for _, v := range s {
		set.Add(v)
	}

	return set
}
