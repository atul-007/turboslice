package turboslice

import (
	"encoding/binary"
	"math"
	"testing"
)

// Fuzz tests cross-check the dispatched implementations (which include the
// SIMD code path on GOEXPERIMENT=simd builds) against the naive scalar
// baselines. The corpus is decoded from arbitrary bytes so the fuzzer can
// freely explore lengths, including the tail boundaries (lane, lane+1,
// lane-1, etc.) where the SIMD vs scalar split happens.
//
// Run with:
//
//	go test -fuzz=FuzzSumInt32 -fuzztime=10s
//	GOEXPERIMENT=simd go test -fuzz=FuzzSumInt32 -fuzztime=10s

func bytesToInt32(b []byte) []int32 {
	n := len(b) / 4
	out := make([]int32, n)
	for i := 0; i < n; i++ {
		out[i] = int32(binary.LittleEndian.Uint32(b[i*4:]))
	}
	return out
}

func bytesToFloat64(b []byte) []float64 {
	n := len(b) / 8
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		bits := binary.LittleEndian.Uint64(b[i*8:])
		out[i] = math.Float64frombits(bits)
	}
	return out
}

func FuzzSumInt32(f *testing.F) {
	f.Add([]byte{0, 0, 0, 1, 0, 0, 0, 2})
	f.Add(make([]byte, 17*4+1))
	f.Fuzz(func(t *testing.T, b []byte) {
		s := bytesToInt32(b)
		want := naiveSumInt32(s)
		if got := Sum(s); got != want {
			t.Fatalf("Sum disagrees: got %d, want %d (n=%d)", got, want, len(s))
		}
		if got := SumInt32(s); got != want {
			t.Fatalf("SumInt32 disagrees: got %d, want %d (n=%d)", got, want, len(s))
		}
	})
}

func FuzzMinInt32(f *testing.F) {
	f.Add([]byte{0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 3})
	f.Fuzz(func(t *testing.T, b []byte) {
		s := bytesToInt32(b)
		if len(s) == 0 {
			return
		}
		want := naiveMinInt32(s)
		if got := Min(s); got != want {
			t.Fatalf("Min disagrees: got %d, want %d (n=%d)", got, want, len(s))
		}
		if got := MinInt32(s); got != want {
			t.Fatalf("MinInt32 disagrees: got %d, want %d (n=%d)", got, want, len(s))
		}
	})
}

func FuzzFindInt32(f *testing.F) {
	f.Add([]byte{0, 0, 0, 5, 0, 0, 0, 5, 0, 0, 0, 5}, int32(5))
	f.Add([]byte{0, 0, 0, 1, 0, 0, 0, 2}, int32(9))
	f.Fuzz(func(t *testing.T, b []byte, val int32) {
		s := bytesToInt32(b)
		want := naiveFindInt32(s, val)
		if got := Find(s, val); got != want {
			t.Fatalf("Find disagrees: got %d, want %d (n=%d val=%d)",
				got, want, len(s), val)
		}
		if got := FindInt32(s, val); got != want {
			t.Fatalf("FindInt32 disagrees: got %d, want %d (n=%d val=%d)",
				got, want, len(s), val)
		}
	})
}

func FuzzCountInt32(f *testing.F) {
	f.Add([]byte{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 2}, int32(1))
	f.Fuzz(func(t *testing.T, b []byte, val int32) {
		s := bytesToInt32(b)
		want := naiveCountInt32(s, val)
		if got := Count(s, val); got != want {
			t.Fatalf("Count disagrees: got %d, want %d (n=%d val=%d)",
				got, want, len(s), val)
		}
	})
}

func FuzzDotProductFloat64(f *testing.F) {
	// Seed: two pairs of float64 (16 bytes each), reused for s1 and s2.
	f.Add(make([]byte, 32))
	f.Fuzz(func(t *testing.T, b []byte) {
		s := bytesToFloat64(b)
		if len(s) < 2 {
			return
		}
		half := len(s) / 2
		a, c := s[:half], s[half:half*2]
		// Skip NaN-containing inputs: float dot ordering differs between
		// scalar-linear and SIMD-tree reductions when NaNs poison the
		// accumulator partway through.
		if containsNonFinite(a) || containsNonFinite(c) {
			return
		}
		want := naiveDotProductFloat64(a, c)
		got := DotProduct(a, c)
		if !approxEqual(got, want, 1e-6) {
			t.Fatalf("DotProduct disagrees: got %g, want %g (n=%d)",
				got, want, len(a))
		}
	})
}

func FuzzAddSlicesInt32(f *testing.F) {
	f.Add(make([]byte, 4*7), make([]byte, 4*7))
	f.Fuzz(func(t *testing.T, b1, b2 []byte) {
		s1 := bytesToInt32(b1)
		s2 := bytesToInt32(b2)
		want := naiveAddSlicesInt32(
			s1[:minLen(s1, s2)],
			s2[:minLen(s1, s2)],
		)
		got := AddSlices(s1, s2)
		if len(got) != len(want) {
			t.Fatalf("AddSlices len: got %d, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("AddSlices[%d]: got %d, want %d", i, got[i], want[i])
			}
		}
	})
}

func containsNonFinite(s []float64) bool {
	for _, v := range s {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

func minLen[T any](a, b []T) int {
	if len(a) < len(b) {
		return len(a)
	}
	return len(b)
}
