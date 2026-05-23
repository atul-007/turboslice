package turboslice

import (
	"math"
	"testing"
)

// Tail-handling correctness tests. The SIMD path processes elements in
// lane-width chunks (4 for int32/float32, 2 for int64/float64) and finishes
// the remainder with a scalar loop. These tests sweep slice lengths around
// every lane boundary so off-by-ones in the tail loop are caught regardless
// of which build path is active.

var tailLengths = []int{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
	23, 31, 32, 33, 63, 64, 65, 127, 128, 129, 255, 257,
	1023, 1024, 1025,
}

// --- Sum tail ---

func TestSumInt32_Tails(t *testing.T) {
	for _, n := range tailLengths {
		s := make([]int32, n)
		var want int32
		for i := range s {
			s[i] = int32(i*37 - 11)
			want += s[i]
		}
		if got := Sum(s); got != want {
			t.Errorf("Sum int32 n=%d: got %d, want %d", n, got, want)
		}
		if got := SumInt32(s); got != want {
			t.Errorf("SumInt32 n=%d: got %d, want %d", n, got, want)
		}
	}
}

func TestSumInt64_Tails(t *testing.T) {
	for _, n := range tailLengths {
		s := make([]int64, n)
		var want int64
		for i := range s {
			s[i] = int64(i)*int64(math.MaxInt32) - 7
			want += s[i]
		}
		if got := Sum(s); got != want {
			t.Errorf("Sum int64 n=%d: got %d, want %d", n, got, want)
		}
	}
}

func TestSumFloat64_Tails(t *testing.T) {
	for _, n := range tailLengths {
		s := make([]float64, n)
		var want float64
		for i := range s {
			s[i] = float64(i) * 0.5
			want += s[i]
		}
		got := Sum(s)
		// Floating-point sum order may differ between SIMD-tree and scalar-linear
		// reductions; allow a small relative tolerance.
		if !approxEqual(got, want, 1e-9) {
			t.Errorf("Sum float64 n=%d: got %g, want %g", n, got, want)
		}
	}
}

// --- Min / Max / MinMax tail ---

func TestMinMaxInt32_Tails(t *testing.T) {
	for _, n := range tailLengths {
		s := make([]int32, n)
		s[0] = 1000
		for i := 1; i < n; i++ {
			s[i] = int32((i * 17) % 73)
		}
		// Stick a min at a tail position and a max somewhere else when possible.
		if n > 4 {
			s[n-1] = -9999
		}
		if n > 6 {
			s[n/2] = 9999
		}
		wantMin, wantMax := naiveMinMax(s)
		if got := Min(s); got != wantMin {
			t.Errorf("Min int32 n=%d: got %d, want %d", n, got, wantMin)
		}
		if got := Max(s); got != wantMax {
			t.Errorf("Max int32 n=%d: got %d, want %d", n, got, wantMax)
		}
		gotLo, gotHi := MinMax(s)
		if gotLo != wantMin || gotHi != wantMax {
			t.Errorf("MinMax int32 n=%d: got (%d,%d), want (%d,%d)",
				n, gotLo, gotHi, wantMin, wantMax)
		}
	}
}

func TestMinMaxFloat64_Tails(t *testing.T) {
	for _, n := range tailLengths {
		s := make([]float64, n)
		for i := range s {
			s[i] = float64((i*31)%97) + 0.25
		}
		if n > 2 {
			s[n-1] = -3.14
		}
		if n > 5 {
			s[n/3] = 999.999
		}
		wantLo := s[0]
		wantHi := s[0]
		for _, v := range s[1:] {
			if v < wantLo {
				wantLo = v
			}
			if v > wantHi {
				wantHi = v
			}
		}
		if got := Min(s); got != wantLo {
			t.Errorf("Min float64 n=%d: got %g, want %g", n, got, wantLo)
		}
		if got := Max(s); got != wantHi {
			t.Errorf("Max float64 n=%d: got %g, want %g", n, got, wantHi)
		}
		gotLo, gotHi := MinMax(s)
		if gotLo != wantLo || gotHi != wantHi {
			t.Errorf("MinMax float64 n=%d: got (%g,%g), want (%g,%g)",
				n, gotLo, gotHi, wantLo, wantHi)
		}
	}
}

// --- Count / Find tail ---

func TestCountInt32_Tails(t *testing.T) {
	for _, n := range tailLengths {
		s := make([]int32, n)
		for i := range s {
			if i%3 == 0 {
				s[i] = 42
			} else {
				s[i] = int32(i)
			}
		}
		want := 0
		for _, v := range s {
			if v == 42 {
				want++
			}
		}
		if got := Count(s, int32(42)); got != want {
			t.Errorf("Count int32 n=%d: got %d, want %d", n, got, want)
		}
	}
}

func TestFindInt32_TailPositions(t *testing.T) {
	// For each tail length, plant the needle at every possible position
	// (including the very last element, which lives in the scalar tail
	// for non-lane-multiple lengths) and confirm Find returns that index.
	for _, n := range tailLengths {
		for pos := 0; pos < n; pos++ {
			s := make([]int32, n)
			for i := range s {
				s[i] = -1
			}
			s[pos] = 42
			if got := Find(s, int32(42)); got != pos {
				t.Errorf("Find int32 n=%d pos=%d: got %d", n, pos, got)
			}
		}
	}
}

// --- AddSlices / MulSlices / DotProduct tail ---

func TestAddSlicesInt32_Tails(t *testing.T) {
	for _, n := range tailLengths {
		a := make([]int32, n)
		b := make([]int32, n)
		want := make([]int32, n)
		for i := range a {
			a[i] = int32(i)
			b[i] = int32(2*i + 1)
			want[i] = a[i] + b[i]
		}
		got := AddSlices(a, b)
		assertSliceEqual(t, got, want)
	}
}

func TestMulSlicesFloat64_Tails(t *testing.T) {
	for _, n := range tailLengths {
		a := make([]float64, n)
		b := make([]float64, n)
		want := make([]float64, n)
		for i := range a {
			a[i] = float64(i) * 0.5
			b[i] = float64(i)*0.25 + 1
			want[i] = a[i] * b[i]
		}
		got := MulSlices(a, b)
		for i := range got {
			if !approxEqual(got[i], want[i], 1e-9) {
				t.Errorf("MulSlices float64 n=%d i=%d: got %g, want %g",
					n, i, got[i], want[i])
				break
			}
		}
	}
}

func TestDotProductInt32_Tails(t *testing.T) {
	for _, n := range tailLengths {
		a := make([]int32, n)
		b := make([]int32, n)
		var want int32
		for i := range a {
			a[i] = int32(i % 11)
			b[i] = int32((i + 3) % 7)
			want += a[i] * b[i]
		}
		if got := DotProduct(a, b); got != want {
			t.Errorf("DotProduct int32 n=%d: got %d, want %d", n, got, want)
		}
	}
}

func TestDotProductFloat64_Tails(t *testing.T) {
	for _, n := range tailLengths {
		a := make([]float64, n)
		b := make([]float64, n)
		var want float64
		for i := range a {
			a[i] = float64(i) * 0.125
			b[i] = float64(i)*0.375 - 1
			want += a[i] * b[i]
		}
		got := DotProduct(a, b)
		if !approxEqual(got, want, 1e-6) {
			t.Errorf("DotProduct float64 n=%d: got %g, want %g", n, got, want)
		}
	}
}

// --- Int64 type coverage (its dot/mul paths are scalar even on SIMD builds) ---

func TestDotProductInt64(t *testing.T) {
	a := []int64{1, 2, 3, 4, 5, 6, 7}
	b := []int64{10, 20, 30, 40, 50, 60, 70}
	want := int64(0)
	for i := range a {
		want += a[i] * b[i]
	}
	if got := DotProduct(a, b); got != want {
		t.Errorf("DotProduct int64 = %d, want %d", got, want)
	}
	if got := DotProductInt64(a, b); got != want {
		t.Errorf("DotProductInt64 = %d, want %d", got, want)
	}
}

func TestMulSlicesInt64(t *testing.T) {
	a := []int64{2, 3, 4, 5, 6}
	b := []int64{10, 10, 10, 10, 10}
	want := []int64{20, 30, 40, 50, 60}
	got := MulSlices(a, b)
	assertSliceEqual(t, got, want)
}

func TestAddSlicesInt64(t *testing.T) {
	a := []int64{math.MaxInt32, math.MaxInt32, math.MaxInt32}
	b := []int64{1, 2, 3}
	want := []int64{math.MaxInt32 + 1, math.MaxInt32 + 2, math.MaxInt32 + 3}
	got := AddSlices(a, b)
	assertSliceEqual(t, got, want)
}

// --- NaN / Inf semantics ---
//
// Sum is required to propagate NaN and to overflow to +/-Inf consistently
// regardless of SIMD vs scalar build. Min/Max with NaN have well-known
// implementation-defined behavior on SSE (MINPS/MAXPS return the second
// operand on NaN); callers should scrub NaNs first and we don't test it
// here on purpose.

func TestSumFloat64_NaNPropagates(t *testing.T) {
	s := []float64{1, 2, math.NaN(), 4, 5, 6, 7, 8, 9}
	if got := Sum(s); !math.IsNaN(got) {
		t.Errorf("Sum with NaN = %g, want NaN", got)
	}
}

func TestSumFloat64_PositiveInfinity(t *testing.T) {
	s := []float64{1, 2, math.Inf(1), 4}
	if got := Sum(s); !math.IsInf(got, 1) {
		t.Errorf("Sum with +Inf = %g, want +Inf", got)
	}
}

func TestSumFloat32_NaNPropagates(t *testing.T) {
	s := []float32{1, 2, float32(math.NaN()), 4, 5}
	if got := Sum(s); !math.IsNaN(float64(got)) {
		t.Errorf("Sum float32 with NaN = %g, want NaN", got)
	}
}

// --- Generic-fallback equivalence ---
//
// Generic Sum/Min/Max on non-SIMD types must agree with the typed naive
// loop. This exercises the type-switch default branch, which previously
// re-asserted any().([]T) inside the loop in AddSlices/MulSlices.

func TestGenericFallback_Int8(t *testing.T) {
	s := []int8{1, -2, 3, -4, 5, -6, 7, -8, 9, -10}
	if got, want := Sum(s), int8(-5); got != want {
		t.Errorf("Sum int8 = %d, want %d", got, want)
	}
	if got, want := Min(s), int8(-10); got != want {
		t.Errorf("Min int8 = %d, want %d", got, want)
	}
}

func TestGenericFallback_AddSlicesInt(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{10, 20, 30, 40, 50}
	want := []int{11, 22, 33, 44, 55}
	got := AddSlices(a, b)
	assertSliceEqual(t, got, want)
}

func TestGenericFallback_MulSlicesInt(t *testing.T) {
	a := []int{2, 3, 4}
	b := []int{5, 6, 7}
	want := []int{10, 18, 28}
	got := MulSlices(a, b)
	assertSliceEqual(t, got, want)
}

// --- Length-mismatch truncation ---

func TestAddSlices_MismatchTruncates(t *testing.T) {
	a := []int32{1, 2, 3, 4, 5, 6}
	b := []int32{10, 20, 30}
	got := AddSlices(a, b)
	want := []int32{11, 22, 33}
	assertSliceEqual(t, got, want)
}

func TestDotProduct_MismatchTruncates(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5}
	b := []float64{1, 1, 1}
	if got, want := DotProduct(a, b), 6.0; got != want {
		t.Errorf("DotProduct mismatch = %g, want %g", got, want)
	}
}

// --- generic.go contract tests ---

func TestFilter_NonNilNonMatch(t *testing.T) {
	// Documented contract: non-nil input always returns non-nil.
	got := Filter([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 })
	if got == nil {
		t.Error("Filter with no matches on non-nil input should return non-nil")
	}
	if len(got) != 0 {
		t.Errorf("Filter with no matches: len=%d, want 0", len(got))
	}
}

func TestFlatten_NilInputReturnsNil(t *testing.T) {
	if Flatten[int](nil) != nil {
		t.Error("Flatten(nil) should return nil")
	}
}

func TestFlatten_NonNilEmptyReturnsNonNil(t *testing.T) {
	got := Flatten([][]int{})
	if got == nil {
		t.Error("Flatten([][]int{}) should return non-nil")
	}
}

// --- helpers ---

func naiveMinMax(s []int32) (int32, int32) {
	lo, hi := s[0], s[0]
	for _, v := range s[1:] {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	return lo, hi
}

func approxEqual(a, b, tol float64) bool {
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	if diff <= tol {
		return true
	}
	scale := math.Abs(a) + math.Abs(b)
	return diff/scale <= tol
}
