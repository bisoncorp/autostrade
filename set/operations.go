package set

import (
	"bytes"
	"fmt"
)

func Union[T comparable](a, b Set[T]) Set[T] {
	r := New[T]()
	a.Filter(func(value T) bool {
		r.Insert(value)
		return false
	})
	b.Filter(func(value T) bool {
		r.Insert(value)
		return false
	})
	return r
}

func Intersection[T comparable](a, b Set[T]) Set[T] {
	r := New[T]()
	a.Filter(func(value T) bool {
		if b.Has(value) {
			r.Insert(value)
		}
		return false
	})
	return r
}

func Subtraction[T comparable](a, b Set[T]) Set[T] {
	r := New[T]()
	a.Filter(func(value T) bool {
		if !b.Has(value) {
			r.Insert(value)
		}
		return false
	})
	return r
}

func ExclusiveOr[T comparable](a, b Set[T]) Set[T] {
	return Subtraction(Union(a, b), Intersection(a, b))
}

func String[T comparable](a Set[T]) string {
	buf := bytes.Buffer{}
	a.Filter(func(value T) bool {
		buf.WriteString(fmt.Sprintf("%v, ", value))
		return false
	})
	str := buf.String()
	str = str[:len(str)-2]
	return fmt.Sprintf("{ %s }", str)
}
