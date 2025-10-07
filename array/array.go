package array

// Filter returns a new slice containing only the elements for which the
// predicate returns true. It preserves the original order and does not modify
// the input slice. Works with any element type (structs, strings, ints, etc.).
//
// If predicate is nil, a shallow copy of the input slice is returned.
func Filter[T any](in []T, predicate func(T) bool) []T {
	if predicate == nil {
		out := make([]T, len(in))
		copy(out, in)
		return out
	}

	out := make([]T, 0, len(in))
	for _, v := range in {
		if predicate(v) {
			out = append(out, v)
		}
	}
	return out
}

// FilterIndexed behaves like Filter but the predicate additionally receives the
// element's index and the original input slice, similar to JavaScript's
// Array.filter callback signature: (item, index, array) => boolean.
//
// If predicate is nil, a shallow copy of the input slice is returned.
func FilterIndexed[T any](in []T, predicate func(v T, index int, arr []T) bool) []T {
	if predicate == nil {
		out := make([]T, len(in))
		copy(out, in)
		return out
	}

	out := make([]T, 0, len(in))
	for i, v := range in {
		if predicate(v, i, in) {
			out = append(out, v)
		}
	}
	return out
}

// Map transforms each element in the input slice using the provided mapper function
// and returns a new slice with the transformed elements. Similar to JavaScript's
// Array.map(). The input and output types can be different.
//
// If mapper is nil, panics (mapper is required for transformation).
func Map[T any, R any](in []T, mapper func(T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = mapper(v)
	}
	return out
}

// MapIndexed behaves like Map but the mapper additionally receives the element's
// index and the original input slice, similar to JavaScript's Array.map callback
// signature: (item, index, array) => newValue.
//
// If mapper is nil, panics (mapper is required for transformation).
func MapIndexed[T any, R any](in []T, mapper func(v T, index int, arr []T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = mapper(v, i, in)
	}
	return out
}

// Includes checks if at least one element in the slice satisfies the provided
// predicate function. Returns true if the predicate returns true for any element,
// false otherwise. Similar to JavaScript's Array.some().
//
// If predicate is nil, returns false.
func Includes[T any](in []T, predicate func(T) bool) bool {
	if predicate == nil {
		return false
	}

	for _, v := range in {
		if predicate(v) {
			return true
		}
	}
	return false
}

// IncludesIndexed behaves like Includes but the predicate additionally receives
// the element's index and the original input slice, similar to JavaScript's
// Array.some callback signature: (item, index, array) => boolean.
//
// If predicate is nil, returns false.
func IncludesIndexed[T any](in []T, predicate func(v T, index int, arr []T) bool) bool {
	if predicate == nil {
		return false
	}

	for i, v := range in {
		if predicate(v, i, in) {
			return true
		}
	}
	return false
}
