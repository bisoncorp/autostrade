package set

type OperationFunc[T comparable] func(value T) T
type FilterFunc[T comparable] func(value T) bool

type Set[T comparable] interface {
	Insert(T)
	Remove(T)
	Has(T) bool
	Len() int
	Do(OperationFunc[T]) Set[T]
	Filter(FilterFunc[T]) Set[T]
}

type empty struct{}
type mapset[T comparable] map[T]empty

func New[T comparable]() Set[T] {
	return mapset[T](make(map[T]empty))
}

func (s mapset[T]) Insert(value T) {
	s[value] = empty{}
}

func (s mapset[T]) Remove(value T) {
	delete(s, value)
}

func (s mapset[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}

func (s mapset[T]) Do(op OperationFunc[T]) Set[T] {
	r := New[T]()
	for k := range s {
		r.Insert(op(k))
	}
	return r
}

func (s mapset[T]) Filter(op FilterFunc[T]) Set[T] {
	r := New[T]()
	for k := range s {
		if op(k) {
			r.Insert(k)
		}
	}
	return r
}

func (s mapset[T]) Len() int {
	return len(s)
}
