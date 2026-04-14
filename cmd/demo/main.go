package main

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"time"

	"github.com/atul-007/turboslice"
)

func main() {
	fmt.Println("=== TurboSlice Performance Demo ===")
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("SIMD:     %s\n", simdStatus())
	fmt.Println()
	fmt.Println("Three variants compared:")
	fmt.Println("  Typed   = turboslice.SumInt32()       (direct, fully inlined)")
	fmt.Println("  Generic = turboslice.Sum[int32]()     (generic dispatch)")
	fmt.Println("  Naive   = hand-written for-loop")
	fmt.Println()

	sizes := []int{1_000, 100_000, 1_000_000, 10_000_000}

	for _, n := range sizes {
		fmt.Printf("--- %s elements ---\n", formatNum(n))
		i32 := randInt32(n)
		f64a := randFloat64(n)
		f64b := randFloat64(n)

		compare3("Sum[int32]",
			func() { turboslice.SumInt32(i32) },
			func() { turboslice.Sum(i32) },
			func() { naiveSum(i32) },
		)
		compare3("Min[int32]",
			func() { turboslice.MinInt32(i32) },
			func() { turboslice.Min(i32) },
			func() { naiveMin(i32) },
		)
		compare3("Max[int32]",
			func() { turboslice.MaxInt32(i32) },
			func() { turboslice.Max(i32) },
			func() { naiveMax(i32) },
		)
		compare3("Find[int32] miss",
			func() { turboslice.FindInt32(i32, -99999) },
			func() { turboslice.Find(i32, int32(-99999)) },
			func() { naiveFind(i32, int32(-99999)) },
		)
		compare3("Count[int32]",
			func() { turboslice.CountInt32(i32, i32[n/2]) },
			func() { turboslice.Count(i32, i32[n/2]) },
			func() { naiveCount(i32, i32[n/2]) },
		)
		compare3("DotProduct[f64]",
			func() { turboslice.DotProductFloat64(f64a, f64b) },
			func() { turboslice.DotProduct(f64a, f64b) },
			func() { naiveDotF64(f64a, f64b) },
		)
		fmt.Println()
	}
}

func compare3(name string, typed, generic, naive func()) {
	tTyped := bench(typed)
	tGeneric := bench(generic)
	tNaive := bench(naive)
	speedupTyped := float64(tNaive) / float64(tTyped)
	speedupGeneric := float64(tNaive) / float64(tGeneric)
	fmt.Printf("  %-22s typed=%-11s (%.2fx)   generic=%-11s (%.2fx)   naive=%-11s\n",
		name, tTyped, speedupTyped, tGeneric, speedupGeneric, tNaive)
}

func bench(fn func()) time.Duration {
	for range 10 {
		fn()
	}
	runtime.GC()

	const iters = 300
	start := time.Now()
	for range iters {
		fn()
	}
	return time.Since(start) / time.Duration(iters)
}

func simdStatus() string {
	if runtime.GOARCH == "amd64" {
		return "available (AMD64)"
	}
	return "off (scalar fallback, " + runtime.GOARCH + ")"
}

func formatNum(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%dM", n/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%dK", n/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

func randInt32(n int) []int32 {
	s := make([]int32, n)
	for i := range s {
		s[i] = rand.Int32N(10000) - 5000
	}
	return s
}

func randFloat64(n int) []float64 {
	s := make([]float64, n)
	for i := range s {
		s[i] = rand.Float64()*1000 - 500
	}
	return s
}

// --- naive baselines ---

func naiveSum(s []int32) int32 {
	var t int32
	for _, v := range s {
		t += v
	}
	return t
}

func naiveFind(s []int32, val int32) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func naiveMin(s []int32) int32 {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func naiveMax(s []int32) int32 {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func naiveCount(s []int32, val int32) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

func naiveDotF64(a, b []float64) float64 {
	var t float64
	for i := range a {
		t += a[i] * b[i]
	}
	return t
}
