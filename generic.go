package turboslice

// Higher-order slice operations. These are pure Go generics —
// no SIMD acceleration — but provide a complete "Standard Library+"
// experience alongside the accelerated numeric functions.

// Map applies fn to each element of s and returns the results.
func Map[T any, U any](s []T, fn func(T) U) []U {
	if s == nil {
		return nil
	}
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

// Filter returns a new slice containing only elements for which fn returns true.
// Contract: nil input returns nil; any non-nil input (even all-rejected) returns
// a non-nil slice. The result is preallocated to len(s) capacity, which is
// optimal for selective filters and wastes at most one allocation otherwise.
func Filter[T any](s []T, fn func(T) bool) []T {
	if s == nil {
		return nil
	}
	result := make([]T, 0, len(s))
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce applies fn to an accumulator and each element (left to right),
// reducing the slice to a single value. The initial value seeds the accumulator.
func Reduce[T any, U any](s []T, initial U, fn func(U, T) U) U {
	acc := initial
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

// Any reports whether fn returns true for any element in s.
func Any[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

// All reports whether fn returns true for every element in s.
// Returns true for empty slices.
func All[T any](s []T, fn func(T) bool) bool {
	for _, v := range s {
		if !fn(v) {
			return false
		}
	}
	return true
}

// ForEach calls fn on each element of s.
func ForEach[T any](s []T, fn func(T)) {
	for _, v := range s {
		fn(v)
	}
}

// Chunk splits s into sub-slices of size n.
// The last chunk may have fewer than n elements.
func Chunk[T any](s []T, n int) [][]T {
	if n <= 0 {
		panic("turboslice: Chunk size must be positive")
	}
	if len(s) == 0 {
		return nil
	}
	chunks := make([][]T, 0, (len(s)+n-1)/n)
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

// Unique returns a new slice with duplicate elements removed,
// preserving the order of first occurrence.
func Unique[T comparable](s []T) []T {
	if s == nil {
		return nil
	}
	seen := make(map[T]struct{}, len(s))
	result := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Reverse returns a new slice with elements in reverse order.
func Reverse[T any](s []T) []T {
	if s == nil {
		return nil
	}
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// Flatten concatenates a slice of slices into a single slice.
// Contract: nil input returns nil; non-nil input always returns non-nil.
func Flatten[T any](ss [][]T) []T {
	if ss == nil {
		return nil
	}
	total := 0
	for _, s := range ss {
		total += len(s)
	}
	result := make([]T, 0, total)
	for _, s := range ss {
		result = append(result, s...)
	}
	return result
}
