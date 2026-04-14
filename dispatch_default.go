//go:build !goexperiment.simd || !amd64

package turboslice

// Fallback implementations for platforms without SIMD support.
// These are optimized scalar loops that the Go compiler can auto-vectorize.

// --- Sum ---

func sumInt32(s []int32) int32 {
	var total int32
	for _, v := range s {
		total += v
	}
	return total
}

func sumInt64(s []int64) int64 {
	var total int64
	for _, v := range s {
		total += v
	}
	return total
}

func sumFloat32(s []float32) float32 {
	var total float32
	for _, v := range s {
		total += v
	}
	return total
}

func sumFloat64(s []float64) float64 {
	var total float64
	for _, v := range s {
		total += v
	}
	return total
}

// --- Find ---

func findInt32(s []int32, val int32) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func findInt64(s []int64, val int64) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func findFloat32(s []float32, val float32) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func findFloat64(s []float64, val float64) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

// --- Count ---

func countInt32(s []int32, val int32) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

func countInt64(s []int64, val int64) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

func countFloat32(s []float32, val float32) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

func countFloat64(s []float64, val float64) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

// --- Min ---

func minInt32(s []int32) int32 {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func minInt64(s []int64) int64 {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func minFloat32(s []float32) float32 {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func minFloat64(s []float64) float64 {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// --- Max ---

func maxInt32(s []int32) int32 {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func maxInt64(s []int64) int64 {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func maxFloat32(s []float32) float32 {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func maxFloat64(s []float64) float64 {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// --- MinMax ---

func minmaxInt32(s []int32) (int32, int32) {
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

func minmaxInt64(s []int64) (int64, int64) {
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

func minmaxFloat32(s []float32) (float32, float32) {
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

func minmaxFloat64(s []float64) (float64, float64) {
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

// --- Element-wise Arithmetic ---

func addSlicesInt32(dst, s1, s2 []int32) {
	for i := range dst {
		dst[i] = s1[i] + s2[i]
	}
}

func addSlicesInt64(dst, s1, s2 []int64) {
	for i := range dst {
		dst[i] = s1[i] + s2[i]
	}
}

func addSlicesFloat32(dst, s1, s2 []float32) {
	for i := range dst {
		dst[i] = s1[i] + s2[i]
	}
}

func addSlicesFloat64(dst, s1, s2 []float64) {
	for i := range dst {
		dst[i] = s1[i] + s2[i]
	}
}

func mulSlicesInt32(dst, s1, s2 []int32) {
	for i := range dst {
		dst[i] = s1[i] * s2[i]
	}
}

func mulSlicesFloat32(dst, s1, s2 []float32) {
	for i := range dst {
		dst[i] = s1[i] * s2[i]
	}
}

func mulSlicesFloat64(dst, s1, s2 []float64) {
	for i := range dst {
		dst[i] = s1[i] * s2[i]
	}
}

// --- Dot Product ---

func dotProductInt32(s1, s2 []int32) int32 {
	var total int32
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}

func dotProductInt64(s1, s2 []int64) int64 {
	var total int64
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}

func dotProductFloat32(s1, s2 []float32) float32 {
	var total float32
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}

func dotProductFloat64(s1, s2 []float64) float64 {
	var total float64
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}
