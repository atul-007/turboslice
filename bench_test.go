package turboslice

import (
	"math/rand/v2"
	"testing"
)

// Benchmark sizes to demonstrate scaling behavior
var benchSizes = []struct {
	name string
	n    int
}{
	{"64", 64},
	{"1K", 1_024},
	{"64K", 65_536},
	{"1M", 1_048_576},
}

// --- Data generators ---

func randInt32Slice(n int) []int32 {
	s := make([]int32, n)
	for i := range s {
		s[i] = rand.Int32N(1000) - 500
	}
	return s
}

func randInt64Slice(n int) []int64 {
	s := make([]int64, n)
	for i := range s {
		s[i] = rand.Int64N(1000) - 500
	}
	return s
}

func randFloat32Slice(n int) []float32 {
	s := make([]float32, n)
	for i := range s {
		s[i] = rand.Float32()*1000 - 500
	}
	return s
}

func randFloat64Slice(n int) []float64 {
	s := make([]float64, n)
	for i := range s {
		s[i] = rand.Float64()*1000 - 500
	}
	return s
}

// --- Naive baseline implementations ---

func naiveSum[T Numeric](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

func naiveFind[T comparable](s []T, val T) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func naiveMin[T Numeric](s []T) T {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func naiveMax[T Numeric](s []T) T {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func naiveCount[T comparable](s []T, val T) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

func naiveDotProduct[T Numeric](s1, s2 []T) T {
	var total T
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}

func naiveAddSlices[T Numeric](s1, s2 []T) []T {
	n := len(s1)
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = s1[i] + s2[i]
	}
	return result
}

// ============================================================
// Sum Benchmarks
// ============================================================

func BenchmarkSumInt32(b *testing.B) {
	for _, sz := range benchSizes {
		data := randInt32Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Sum(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveSum(data)
			}
		})
	}
}

func BenchmarkSumFloat64(b *testing.B) {
	for _, sz := range benchSizes {
		data := randFloat64Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Sum(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveSum(data)
			}
		})
	}
}

// ============================================================
// Find Benchmarks
// ============================================================

func BenchmarkFindInt32(b *testing.B) {
	for _, sz := range benchSizes {
		data := randInt32Slice(sz.n)
		// Search for last element (worst case)
		target := data[sz.n-1]
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Find(data, target)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveFind(data, target)
			}
		})
	}
}

func BenchmarkFindInt32_NotFound(b *testing.B) {
	for _, sz := range benchSizes {
		data := make([]int32, sz.n)
		for i := range data {
			data[i] = int32(i)
		}
		target := int32(-1) // not in slice
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Find(data, target)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveFind(data, target)
			}
		})
	}
}

// ============================================================
// Count Benchmarks
// ============================================================

func BenchmarkCountInt32(b *testing.B) {
	for _, sz := range benchSizes {
		data := randInt32Slice(sz.n)
		target := data[0]
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Count(data, target)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveCount(data, target)
			}
		})
	}
}

// ============================================================
// Min/Max Benchmarks
// ============================================================

func BenchmarkMinInt32(b *testing.B) {
	for _, sz := range benchSizes {
		data := randInt32Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Min(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveMin(data)
			}
		})
	}
}

func BenchmarkMaxFloat64(b *testing.B) {
	for _, sz := range benchSizes {
		data := randFloat64Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Max(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveMax(data)
			}
		})
	}
}

func BenchmarkMinMaxInt32(b *testing.B) {
	for _, sz := range benchSizes {
		data := randInt32Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				MinMax(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				lo, hi := data[0], data[0]
				for _, v := range data[1:] {
					if v < lo {
						lo = v
					}
					if v > hi {
						hi = v
					}
				}
				_ = lo
				_ = hi
			}
		})
	}
}

// ============================================================
// DotProduct Benchmarks
// ============================================================

func BenchmarkDotProductFloat64(b *testing.B) {
	for _, sz := range benchSizes {
		a := randFloat64Slice(sz.n)
		c := randFloat64Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				DotProduct(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveDotProduct(a, c)
			}
		})
	}
}

func BenchmarkDotProductInt32(b *testing.B) {
	for _, sz := range benchSizes {
		a := randInt32Slice(sz.n)
		c := randInt32Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				DotProduct(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveDotProduct(a, c)
			}
		})
	}
}

// ============================================================
// AddSlices Benchmarks
// ============================================================

func BenchmarkAddSlicesInt32(b *testing.B) {
	for _, sz := range benchSizes {
		a := randInt32Slice(sz.n)
		c := randInt32Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				AddSlices(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveAddSlices(a, c)
			}
		})
	}
}

func BenchmarkAddSlicesFloat64(b *testing.B) {
	for _, sz := range benchSizes {
		a := randFloat64Slice(sz.n)
		c := randFloat64Slice(sz.n)
		b.Run("TurboSlice/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				AddSlices(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveAddSlices(a, c)
			}
		})
	}
}
