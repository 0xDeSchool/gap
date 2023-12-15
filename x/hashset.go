package x

// Set is a set of items.
// It is implemented as a map[T]struct{}.
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	set := make(Set[T])
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

func (set Set[T]) Set(item T) {
	set[item] = struct{}{}
}

func (set Set[T]) Delete(item T) {
	delete(set, item)
}

func (set Set[T]) Contains(item T) bool {
	_, ok := set[item]
	return ok
}

func (set Set[T]) Len() int {
	return len(set)
}

func (set Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(set))
	for item := range set {
		slice = append(slice, item)
	}
	return slice
}
