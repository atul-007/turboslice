package turboslice

import (
	"math"
	"testing"
)

// --- Sum Tests ---

func TestSumInt32(t *testing.T) {
	tests := []struct {
		name string
		in   []int32
		want int32
	}{
		{"empty", nil, 0},
		{"single", []int32{42}, 42},
		{"small", []int32{1, 2, 3, 4, 5}, 15},
		{"negative", []int32{-1, -2, -3}, -6},
		{"mixed", []int32{-10, 20, -30, 40}, 20},
		{"zeros", []int32{0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{"exactly4", []int32{1, 2, 3, 4}, 10},
		{"exactly8", []int32{1, 2, 3, 4, 5, 6, 7, 8}, 36},
		{"7elements", []int32{1, 2, 3, 4, 5, 6, 7}, 28},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.in)
			if got != tt.want {
				t.Errorf("Sum(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestSumInt64(t *testing.T) {
	tests := []struct {
		name string
		in   []int64
		want int64
	}{
		{"empty", nil, 0},
		{"large_values", []int64{math.MaxInt32 + 1, math.MaxInt32 + 1}, 2 * (math.MaxInt32 + 1)},
		{"mixed", []int64{-100, 200, -300, 400, -500}, -300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.in)
			if got != tt.want {
				t.Errorf("Sum(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestSumFloat64(t *testing.T) {
	tests := []struct {
		name string
		in   []float64
		want float64
	}{
		{"empty", nil, 0},
		{"simple", []float64{1.5, 2.5, 3.0}, 7.0},
		{"negative", []float64{-1.1, -2.2, -3.3}, -6.6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.in)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Sum(%v) = %f, want %f", tt.in, got, tt.want)
			}
		})
	}
}

func TestSumFloat32(t *testing.T) {
	got := Sum([]float32{1.0, 2.0, 3.0, 4.0, 5.0})
	if got != 15.0 {
		t.Errorf("Sum(float32) = %f, want 15.0", got)
	}
}

func TestSumOtherTypes(t *testing.T) {
	// Tests the generic scalar fallback path
	t.Run("int", func(t *testing.T) {
		got := Sum([]int{1, 2, 3})
		if got != 6 {
			t.Errorf("Sum([]int) = %d, want 6", got)
		}
	})
	t.Run("uint32", func(t *testing.T) {
		got := Sum([]uint32{10, 20, 30})
		if got != 60 {
			t.Errorf("Sum([]uint32) = %d, want 60", got)
		}
	})
	t.Run("int8", func(t *testing.T) {
		got := Sum([]int8{1, 2, 3, 4})
		if got != 10 {
			t.Errorf("Sum([]int8) = %d, want 10", got)
		}
	})
}

func TestSumLargeSlice(t *testing.T) {
	n := 10_000
	s := make([]int32, n)
	for i := range s {
		s[i] = 1
	}
	got := Sum(s)
	if got != int32(n) {
		t.Errorf("Sum(10000 ones) = %d, want %d", got, n)
	}
}

// --- Find Tests ---

func TestFindInt32(t *testing.T) {
	tests := []struct {
		name string
		in   []int32
		val  int32
		want int
	}{
		{"empty", nil, 42, -1},
		{"found_first", []int32{1, 2, 3}, 1, 0},
		{"found_middle", []int32{1, 2, 3}, 2, 1},
		{"found_last", []int32{1, 2, 3}, 3, 2},
		{"not_found", []int32{1, 2, 3}, 4, -1},
		{"duplicate", []int32{1, 2, 2, 3}, 2, 1},
		{"large_aligned", make8Int32(100, 42, 50), 42, 50},
		{"large_not_found", make8Int32(100, 0, -1), 42, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Find(tt.in, tt.val)
			if got != tt.want {
				t.Errorf("Find(%v, %d) = %d, want %d", tt.in, tt.val, got, tt.want)
			}
		})
	}
}

func TestFindInt64(t *testing.T) {
	got := Find([]int64{10, 20, 30, 40, 50}, int64(30))
	if got != 2 {
		t.Errorf("Find(int64, 30) = %d, want 2", got)
	}
}

func TestFindFloat64(t *testing.T) {
	got := Find([]float64{1.1, 2.2, 3.3}, 2.2)
	if got != 1 {
		t.Errorf("Find(float64, 2.2) = %d, want 1", got)
	}
}

func TestFindString(t *testing.T) {
	// Tests generic fallback path (non-numeric)
	got := Find([]string{"a", "b", "c"}, "b")
	if got != 1 {
		t.Errorf("Find(string, b) = %d, want 1", got)
	}
}

// --- Contains Tests ---

func TestContains(t *testing.T) {
	if !Contains([]int32{1, 2, 3}, int32(2)) {
		t.Error("Contains should find 2")
	}
	if Contains([]int32{1, 2, 3}, int32(4)) {
		t.Error("Contains should not find 4")
	}
	if Contains([]int32(nil), int32(1)) {
		t.Error("Contains of nil should be false")
	}
}

// --- Count Tests ---

func TestCountInt32(t *testing.T) {
	tests := []struct {
		name string
		in   []int32
		val  int32
		want int
	}{
		{"empty", nil, 1, 0},
		{"none", []int32{1, 2, 3}, 4, 0},
		{"one", []int32{1, 2, 3}, 2, 1},
		{"many", []int32{5, 5, 5, 5, 5}, 5, 5},
		{"mixed", []int32{1, 2, 1, 2, 1}, 1, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Count(tt.in, tt.val)
			if got != tt.want {
				t.Errorf("Count(%v, %d) = %d, want %d", tt.in, tt.val, got, tt.want)
			}
		})
	}
}

func TestCountLargeSlice(t *testing.T) {
	n := 10_000
	s := make([]int32, n)
	for i := range s {
		if i%3 == 0 {
			s[i] = 42
		}
	}
	expected := (n + 2) / 3
	got := Count(s, int32(42))
	if got != expected {
		t.Errorf("Count(large, 42) = %d, want %d", got, expected)
	}
}

// --- Min/Max Tests ---

func TestMinInt32(t *testing.T) {
	tests := []struct {
		name string
		in   []int32
		want int32
	}{
		{"single", []int32{5}, 5},
		{"sorted", []int32{1, 2, 3, 4}, 1},
		{"reverse", []int32{4, 3, 2, 1}, 1},
		{"negative", []int32{-1, -5, -3}, -5},
		{"mixed", []int32{10, -20, 30, -40, 50}, -40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Min(tt.in)
			if got != tt.want {
				t.Errorf("Min(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestMaxInt32(t *testing.T) {
	tests := []struct {
		name string
		in   []int32
		want int32
	}{
		{"single", []int32{5}, 5},
		{"sorted", []int32{1, 2, 3, 4}, 4},
		{"reverse", []int32{4, 3, 2, 1}, 4},
		{"negative", []int32{-1, -5, -3}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Max(tt.in)
			if got != tt.want {
				t.Errorf("Max(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestMinMaxFloat64(t *testing.T) {
	lo, hi := MinMax([]float64{3.14, -2.71, 1.41, 0.0, 99.9})
	if lo != -2.71 {
		t.Errorf("MinMax min = %f, want -2.71", lo)
	}
	if hi != 99.9 {
		t.Errorf("MinMax max = %f, want 99.9", hi)
	}
}

func TestMinPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Min(nil) should panic")
		}
	}()
	Min([]int32(nil))
}

func TestMaxPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Max(nil) should panic")
		}
	}()
	Max([]int32(nil))
}

func TestMinMaxLargeSlice(t *testing.T) {
	n := 10_000
	s := make([]int32, n)
	for i := range s {
		s[i] = int32(i)
	}
	s[7777] = -999
	s[3333] = 99999

	lo := Min(s)
	hi := Max(s)
	if lo != -999 {
		t.Errorf("Min(large) = %d, want -999", lo)
	}
	if hi != 99999 {
		t.Errorf("Max(large) = %d, want 99999", hi)
	}
}

// --- AddSlices Tests ---

func TestAddSlices(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		a := []int32{1, 2, 3, 4, 5, 6, 7, 8}
		b := []int32{10, 20, 30, 40, 50, 60, 70, 80}
		got := AddSlices(a, b)
		want := []int32{11, 22, 33, 44, 55, 66, 77, 88}
		assertSliceEqual(t, got, want)
	})

	t.Run("float64", func(t *testing.T) {
		a := []float64{1.0, 2.0, 3.0}
		b := []float64{0.5, 0.5, 0.5}
		got := AddSlices(a, b)
		want := []float64{1.5, 2.5, 3.5}
		assertSliceEqual(t, got, want)
	})

	t.Run("different_lengths", func(t *testing.T) {
		a := []int32{1, 2, 3, 4, 5}
		b := []int32{10, 20, 30}
		got := AddSlices(a, b)
		if len(got) != 3 {
			t.Errorf("len = %d, want 3", len(got))
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := AddSlices([]int32{}, []int32{})
		if got != nil {
			t.Error("expected nil for empty input")
		}
	})
}

// --- MulSlices Tests ---

func TestMulSlices(t *testing.T) {
	a := []int32{2, 3, 4, 5}
	b := []int32{10, 10, 10, 10}
	got := MulSlices(a, b)
	want := []int32{20, 30, 40, 50}
	assertSliceEqual(t, got, want)
}

// --- DotProduct Tests ---

func TestDotProduct(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		a := []int32{1, 2, 3, 4}
		b := []int32{5, 6, 7, 8}
		got := DotProduct(a, b)
		want := int32(1*5 + 2*6 + 3*7 + 4*8) // 70
		if got != want {
			t.Errorf("DotProduct = %d, want %d", got, want)
		}
	})

	t.Run("float64", func(t *testing.T) {
		a := []float64{1.0, 2.0, 3.0}
		b := []float64{4.0, 5.0, 6.0}
		got := DotProduct(a, b)
		want := 32.0
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("DotProduct = %f, want %f", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := DotProduct([]int32{}, []int32{})
		if got != 0 {
			t.Errorf("DotProduct(empty) = %d, want 0", got)
		}
	})
}

func TestDotProductLarge(t *testing.T) {
	n := 1000
	a := make([]int32, n)
	b := make([]int32, n)
	for i := range a {
		a[i] = int32(i + 1)
		b[i] = 1
	}
	got := DotProduct(a, b)
	want := int32(n * (n + 1) / 2) // Sum 1..n
	if got != want {
		t.Errorf("DotProduct = %d, want %d", got, want)
	}
}

// --- Generic Utility Tests ---

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(x int) int { return x * 2 })
	assertSliceEqual(t, got, []int{2, 4, 6})

	// Nil input
	if Map[int, int](nil, func(x int) int { return x }) != nil {
		t.Error("Map(nil) should return nil")
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5, 6}, func(x int) bool { return x%2 == 0 })
	assertSliceEqual(t, got, []int{2, 4, 6})

	got2 := Filter([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 })
	if len(got2) != 0 {
		t.Error("Filter should return empty when nothing matches")
	}
}

func TestReduce(t *testing.T) {
	sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, v int) int { return acc + v })
	if sum != 10 {
		t.Errorf("Reduce sum = %d, want 10", sum)
	}

	product := Reduce([]int{1, 2, 3, 4}, 1, func(acc, v int) int { return acc * v })
	if product != 24 {
		t.Errorf("Reduce product = %d, want 24", product)
	}
}

func TestAny(t *testing.T) {
	if !Any([]int{1, 2, 3}, func(x int) bool { return x > 2 }) {
		t.Error("Any should find element > 2")
	}
	if Any([]int{1, 2, 3}, func(x int) bool { return x > 5 }) {
		t.Error("Any should not find element > 5")
	}
}

func TestAll(t *testing.T) {
	if !All([]int{2, 4, 6}, func(x int) bool { return x%2 == 0 }) {
		t.Error("All even should be true")
	}
	if All([]int{2, 3, 6}, func(x int) bool { return x%2 == 0 }) {
		t.Error("All even should be false with 3")
	}
	if !All([]int{}, func(x int) bool { return false }) {
		t.Error("All of empty should be true")
	}
}

func TestChunk(t *testing.T) {
	got := Chunk([]int{1, 2, 3, 4, 5}, 2)
	if len(got) != 3 {
		t.Fatalf("Chunk len = %d, want 3", len(got))
	}
	assertSliceEqual(t, got[0], []int{1, 2})
	assertSliceEqual(t, got[1], []int{3, 4})
	assertSliceEqual(t, got[2], []int{5})
}

func TestChunkPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Chunk(0) should panic")
		}
	}()
	Chunk([]int{1}, 0)
}

func TestUnique(t *testing.T) {
	got := Unique([]int{1, 2, 3, 2, 1, 4})
	assertSliceEqual(t, got, []int{1, 2, 3, 4})
}

func TestReverse(t *testing.T) {
	got := Reverse([]int{1, 2, 3, 4})
	assertSliceEqual(t, got, []int{4, 3, 2, 1})

	if Reverse[int](nil) != nil {
		t.Error("Reverse(nil) should return nil")
	}
}

func TestFlatten(t *testing.T) {
	got := Flatten([][]int{{1, 2}, {3}, {4, 5, 6}})
	assertSliceEqual(t, got, []int{1, 2, 3, 4, 5, 6})
}

// --- Helpers ---

func make8Int32(n int, val int32, pos int) []int32 {
	s := make([]int32, n)
	if pos >= 0 && pos < n {
		s[pos] = val
	}
	return s
}

func assertSliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
