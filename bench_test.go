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

// --- Naive baseline implementations (typed, not generic) ---

func naiveSumInt32(s []int32) int32 {
	var total int32
	for _, v := range s {
		total += v
	}
	return total
}

func naiveFindInt32(s []int32, val int32) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func naiveMinInt32(s []int32) int32 {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func naiveMaxFloat64(s []float64) float64 {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func naiveCountInt32(s []int32, val int32) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

func naiveDotProductFloat64(s1, s2 []float64) float64 {
	var total float64
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}

func naiveDotProductInt32(s1, s2 []int32) int32 {
	var total int32
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}

func naiveAddSlicesInt32(s1, s2 []int32) []int32 {
	n := len(s1)
	result := make([]int32, n)
	for i := 0; i < n; i++ {
		result[i] = s1[i] + s2[i]
	}
	return result
}

func naiveAddSlicesFloat64(s1, s2 []float64) []float64 {
	n := len(s1)
	result := make([]float64, n)
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
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				SumInt32(data)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Sum(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveSumInt32(data)
			}
		})
	}
}

func BenchmarkSumFloat64(b *testing.B) {
	for _, sz := range benchSizes {
		data := randFloat64Slice(sz.n)
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				SumFloat64(data)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Sum(data)
			}
		})
	}
}

// ============================================================
// Find Benchmarks
// ============================================================

func BenchmarkFindInt32_NotFound(b *testing.B) {
	for _, sz := range benchSizes {
		data := make([]int32, sz.n)
		for i := range data {
			data[i] = int32(i)
		}
		target := int32(-1) // not in slice
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				FindInt32(data, target)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Find(data, target)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveFindInt32(data, target)
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
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				CountInt32(data, target)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Count(data, target)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveCountInt32(data, target)
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
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				MinInt32(data)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Min(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveMinInt32(data)
			}
		})
	}
}

func BenchmarkMaxFloat64(b *testing.B) {
	for _, sz := range benchSizes {
		data := randFloat64Slice(sz.n)
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				MaxFloat64(data)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				Max(data)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveMaxFloat64(data)
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
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				DotProductFloat64(a, c)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				DotProduct(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveDotProductFloat64(a, c)
			}
		})
	}
}

func BenchmarkDotProductInt32(b *testing.B) {
	for _, sz := range benchSizes {
		a := randInt32Slice(sz.n)
		c := randInt32Slice(sz.n)
		b.Run("Typed/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				DotProductInt32(a, c)
			}
		})
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				DotProduct(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveDotProductInt32(a, c)
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
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				AddSlices(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveAddSlicesInt32(a, c)
			}
		})
	}
}

func BenchmarkAddSlicesFloat64(b *testing.B) {
	for _, sz := range benchSizes {
		a := randFloat64Slice(sz.n)
		c := randFloat64Slice(sz.n)
		b.Run("Generic/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				AddSlices(a, c)
			}
		})
		b.Run("NaiveLoop/"+sz.name, func(b *testing.B) {
			for b.Loop() {
				naiveAddSlicesFloat64(a, c)
			}
		})
	}
}
