package turboslice

// Scalar generic implementations used as fallback for non-SIMD types
// and as the default path on non-AMD64 architectures.

func sumScalar[T Numeric](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

func findScalar[T comparable](s []T, val T) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}
	return -1
}

func countScalar[T comparable](s []T, val T) int {
	n := 0
	for _, v := range s {
		if v == val {
			n++
		}
	}
	return n
}

func minScalar[T Numeric](s []T) T {
	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func maxScalar[T Numeric](s []T) T {
	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func minmaxScalar[T Numeric](s []T) (T, T) {
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

func dotProductScalar[T Numeric](s1, s2 []T) T {
	var total T
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}
