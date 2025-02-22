package micrograd

type Pair[T any] struct {
	first  T
	second T
}

func NewPair[T any](x, y T) Pair[T] {
	return Pair[T]{first: x, second: y}
}

func (p Pair[T]) X() T {
	return p.first
}

func (p Pair[T]) Y() T {
	return p.second
}

func (p Pair[T]) At(i int) T {
	switch i {
	case 0:
		return p.first
	case 1:
		return p.second
	default:
		panic("index out of range")
	}
}
