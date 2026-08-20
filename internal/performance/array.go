package performance

type FastArray[T any] struct {
	items []T
	chunkSize int
	used int
	capacity int
}

func NewFastArray[T any](chunkSize int) *FastArray[T] {
	return &FastArray[T]{
		items: make([]T, 0, chunkSize),
		chunkSize: chunkSize,
		used: 0,
		capacity: chunkSize,
	}
}

func (fa *FastArray[T]) Append(item T) {
	if fa.used == fa.capacity {
		newCapacity := fa.capacity * 2
		newItems := make([]T, fa.used, newCapacity)
		copy(newItems, fa.items[:fa.used])

		fa.items = newItems[:fa.used]
		fa.capacity = newCapacity
	}

	fa.items = fa.items[:fa.used+1]
	fa.items[fa.used] = item
	fa.used++
}

func (fa *FastArray[T]) Items() []T {
	return fa.items[:fa.used]
}

func (fa *FastArray[T]) Len() int {
	return fa.used
}

func (fa *FastArray[T]) First() T {
	if fa.used == 0 {
		var zero T
		return zero
	}

	return fa.items[0]
}

func (fa *FastArray[T]) Last() T {
	if fa.used == 0 {
		var zero T
		return zero
	}

	return fa.items[fa.used-1]
}
