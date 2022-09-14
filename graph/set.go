package graph

type OrderedSet[T comparable] struct {
	data []T
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		data: make([]T, 0),
	}
}

func (s *OrderedSet[T]) Add(value T) {
	if !s.Has(value) {
		s.data = append(s.data, value)
	}
}

func (s *OrderedSet[T]) Get(index int) T {
	return s.data[index]
}

func (s *OrderedSet[T]) Len() int {
	return len(s.data)
}

func (s *OrderedSet[T]) Delete(index int) {
	s.data = append(s.data[:index], s.data[index+1:]...)
}

func (s *OrderedSet[T]) Clear() {
	s.data = make([]T, 0)
}

func (s *OrderedSet[T]) Has(value T) bool {
	for _, v := range s.data {
		if v == value {
			return true
		}
	}
	return false
}

func (s *OrderedSet[T]) Values() []T {
	return s.data
}

func (s *OrderedSet[T]) Remove(value T) {
	for i, v := range s.data {
		if v == value {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return
		}
	}
}


