// Package turboslice provides SIMD-accelerated slice operations for Go.
//
// TurboSlice is a "Standard Library+" that provides common operations
// (Sum, Find, Filter, Map, Min, Max, DotProduct, etc.) specifically optimized
// using SIMD vector instructions on AMD64 via Go 1.26's experimental
// simd/archsimd package.
//
// On architectures without SIMD support (ARM64, etc.) or when built without
// GOEXPERIMENT=simd, all functions automatically fall back to optimized
// scalar Go code — making this library safe for production on any platform.
//
// Build with SIMD acceleration:
//
//	GOEXPERIMENT=simd go build ./...
//
// Build without (auto-fallback):
//
//	go build ./...
package turboslice

// Numeric is the constraint for all numeric types supported by TurboSlice.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Integer is the constraint for integer types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Float is the constraint for floating-point types.
type Float interface {
	~float32 | ~float64
}

// Sum returns the sum of all elements in the slice.
// Returns the zero value for empty slices.
// For SIMD-accelerated types (int32, int64, float32, float64),
// this uses vectorized addition on AMD64.
func Sum[T Numeric](s []T) T {
	if len(s) == 0 {
		return 0
	}
	switch s := any(s).(type) {
	case []int32:
		return T(sumInt32(s))
	case []int64:
		return T(sumInt64(s))
	case []float32:
		return T(sumFloat32(s))
	case []float64:
		return T(sumFloat64(s))
	default:
		return sumScalar[T](any(s).([]T))
	}
}

// Find returns the index of the first occurrence of val in s,
// or -1 if val is not present.
// For SIMD-accelerated types, this uses vectorized comparison on AMD64.
func Find[T comparable](s []T, val T) int {
	if len(s) == 0 {
		return -1
	}
	switch s := any(s).(type) {
	case []int32:
		return findInt32(s, any(val).(int32))
	case []int64:
		return findInt64(s, any(val).(int64))
	case []float32:
		return findFloat32(s, any(val).(float32))
	case []float64:
		return findFloat64(s, any(val).(float64))
	default:
		return findScalar(any(s).([]T), val)
	}
}

// Contains reports whether val is present in s.
func Contains[T comparable](s []T, val T) bool {
	return Find(s, val) >= 0
}

// Count returns the number of occurrences of val in s.
// For SIMD-accelerated types, this uses vectorized comparison on AMD64.
func Count[T comparable](s []T, val T) int {
	if len(s) == 0 {
		return 0
	}
	switch s := any(s).(type) {
	case []int32:
		return countInt32(s, any(val).(int32))
	case []int64:
		return countInt64(s, any(val).(int64))
	case []float32:
		return countFloat32(s, any(val).(float32))
	case []float64:
		return countFloat64(s, any(val).(float64))
	default:
		return countScalar(any(s).([]T), val)
	}
}

// Min returns the minimum value in the slice.
// Panics if s is empty.
func Min[T Numeric](s []T) T {
	if len(s) == 0 {
		panic("turboslice: Min of empty slice")
	}
	switch s := any(s).(type) {
	case []int32:
		return T(minInt32(s))
	case []int64:
		return T(minInt64(s))
	case []float32:
		return T(minFloat32(s))
	case []float64:
		return T(minFloat64(s))
	default:
		return minScalar[T](any(s).([]T))
	}
}

// Max returns the maximum value in the slice.
// Panics if s is empty.
func Max[T Numeric](s []T) T {
	if len(s) == 0 {
		panic("turboslice: Max of empty slice")
	}
	switch s := any(s).(type) {
	case []int32:
		return T(maxInt32(s))
	case []int64:
		return T(maxInt64(s))
	case []float32:
		return T(maxFloat32(s))
	case []float64:
		return T(maxFloat64(s))
	default:
		return maxScalar[T](any(s).([]T))
	}
}

// MinMax returns both the minimum and maximum values in one pass.
// Panics if s is empty.
func MinMax[T Numeric](s []T) (min, max T) {
	if len(s) == 0 {
		panic("turboslice: MinMax of empty slice")
	}
	switch s := any(s).(type) {
	case []int32:
		lo, hi := minmaxInt32(s)
		return T(lo), T(hi)
	case []int64:
		lo, hi := minmaxInt64(s)
		return T(lo), T(hi)
	case []float32:
		lo, hi := minmaxFloat32(s)
		return T(lo), T(hi)
	case []float64:
		lo, hi := minmaxFloat64(s)
		return T(lo), T(hi)
	default:
		return minmaxScalar[T](any(s).([]T))
	}
}

// AddSlices returns a new slice where each element is s1[i] + s2[i].
// The result length is min(len(s1), len(s2)); extra elements in the longer
// slice are silently ignored. Pass equal-length slices to avoid surprises.
func AddSlices[T Numeric](s1, s2 []T) []T {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}
	if n == 0 {
		return nil
	}
	result := make([]T, n)
	switch v := any(s1).(type) {
	case []int32:
		addSlicesInt32(any(result).([]int32), v[:n], any(s2).([]int32)[:n])
	case []int64:
		addSlicesInt64(any(result).([]int64), v[:n], any(s2).([]int64)[:n])
	case []float32:
		addSlicesFloat32(any(result).([]float32), v[:n], any(s2).([]float32)[:n])
	case []float64:
		addSlicesFloat64(any(result).([]float64), v[:n], any(s2).([]float64)[:n])
	default:
		for i := 0; i < n; i++ {
			result[i] = s1[i] + s2[i]
		}
	}
	return result
}

// MulSlices returns a new slice where each element is s1[i] * s2[i].
// The result length is min(len(s1), len(s2)); extra elements in the longer
// slice are silently ignored. Pass equal-length slices to avoid surprises.
//
// SIMD acceleration covers int32, float32, and float64. int64 multiplication
// has no SSE/AVX2 instruction and uses the scalar path.
func MulSlices[T Numeric](s1, s2 []T) []T {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}
	if n == 0 {
		return nil
	}
	result := make([]T, n)
	switch v := any(s1).(type) {
	case []int32:
		mulSlicesInt32(any(result).([]int32), v[:n], any(s2).([]int32)[:n])
	case []int64:
		mulSlicesInt64(any(result).([]int64), v[:n], any(s2).([]int64)[:n])
	case []float32:
		mulSlicesFloat32(any(result).([]float32), v[:n], any(s2).([]float32)[:n])
	case []float64:
		mulSlicesFloat64(any(result).([]float64), v[:n], any(s2).([]float64)[:n])
	default:
		for i := 0; i < n; i++ {
			result[i] = s1[i] * s2[i]
		}
	}
	return result
}

// DotProduct returns the dot product (sum of element-wise products) of two slices.
// Uses min(len(s1), len(s2)) elements; extra elements in the longer slice are
// silently ignored.
//
// SIMD acceleration covers int32, float32, and float64. int64 dot product
// falls back to a scalar loop (no SSE/AVX2 int64 multiply).
//
// The accumulator type is T, so callers using narrow integer types should
// be wary of overflow on large slices (same behavior as a hand-written loop).
func DotProduct[T Numeric](s1, s2 []T) T {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}
	if n == 0 {
		return 0
	}
	switch s1 := any(s1).(type) {
	case []int32:
		return T(dotProductInt32(s1[:n], any(s2).([]int32)[:n]))
	case []int64:
		return T(dotProductInt64(s1[:n], any(s2).([]int64)[:n]))
	case []float32:
		return T(dotProductFloat32(s1[:n], any(s2).([]float32)[:n]))
	case []float64:
		return T(dotProductFloat64(s1[:n], any(s2).([]float64)[:n]))
	default:
		return dotProductScalar[T](any(s1).([]T)[:n], any(s2).([]T)[:n])
	}
}
