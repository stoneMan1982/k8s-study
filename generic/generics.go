package generic

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// use type parameters to create a generic function
func SumIntsOrFloats[T comparable, V int64 | float64](m map[T]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// user standard library's constraints
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() T {
	n := len(s.data)
	if n == 0 {
		var zero T
		return zero
	}

	v := s.data[n-1]
	s.data = s.data[:n-1]
	return v
}

// interface constraints
type Stringer interface {
	String() string
}

func Print[T Stringer](s T) {
	fmt.Println(s.String())
}

// union types constraints
type Integer interface {
	~int | ~int64
}
