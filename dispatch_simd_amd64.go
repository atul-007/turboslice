//go:build goexperiment.simd && amd64

package turboslice

import (
	"math/bits"
	"simd/archsimd"
)

// SIMD-accelerated implementations using Go 1.26's simd/archsimd package.
// These use 128-bit SSE vectors as the baseline (available on all AMD64),
// with the same API patterns adaptable to 256-bit AVX2 in the future.

const (
	int32Lane  = 4 // Int32x4: 4 elements per 128-bit vector
	int64Lane  = 2 // Int64x2: 2 elements per 128-bit vector
	f32Lane    = 4 // Float32x4: 4 elements per 128-bit vector
	f64Lane    = 2 // Float64x2: 2 elements per 128-bit vector
)

// --- Sum ---

func sumInt32(s []int32) int32 {
	var acc archsimd.Int32x4
	i := 0
	n := len(s)

	// Process 4 elements at a time
	for ; i+int32Lane <= n; i += int32Lane {
		v := archsimd.LoadInt32x4Slice(s[i:])
		acc = acc.Add(v)
	}

	// Horizontal reduction
	var buf [4]int32
	acc.Store(&buf)
	total := buf[0] + buf[1] + buf[2] + buf[3]

	// Handle remaining elements
	for ; i < n; i++ {
		total += s[i]
	}
	return total
}

func sumInt64(s []int64) int64 {
	var acc archsimd.Int64x2
	i := 0
	n := len(s)

	for ; i+int64Lane <= n; i += int64Lane {
		v := archsimd.LoadInt64x2Slice(s[i:])
		acc = acc.Add(v)
	}

	var buf [2]int64
	acc.Store(&buf)
	total := buf[0] + buf[1]

	for ; i < n; i++ {
		total += s[i]
	}
	return total
}

func sumFloat32(s []float32) float32 {
	var acc archsimd.Float32x4
	i := 0
	n := len(s)

	for ; i+f32Lane <= n; i += f32Lane {
		v := archsimd.LoadFloat32x4Slice(s[i:])
		acc = acc.Add(v)
	}

	var buf [4]float32
	acc.Store(&buf)
	total := buf[0] + buf[1] + buf[2] + buf[3]

	for ; i < n; i++ {
		total += s[i]
	}
	return total
}

func sumFloat64(s []float64) float64 {
	var acc archsimd.Float64x2
	i := 0
	n := len(s)

	for ; i+f64Lane <= n; i += f64Lane {
		v := archsimd.LoadFloat64x2Slice(s[i:])
		acc = acc.Add(v)
	}

	var buf [2]float64
	acc.Store(&buf)
	total := buf[0] + buf[1]

	for ; i < n; i++ {
		total += s[i]
	}
	return total
}

// --- Find ---

func findInt32(s []int32, val int32) int {
	n := len(s)
	i := 0

	// Broadcast the search value to all lanes
	var valBuf [4]int32
	valBuf[0], valBuf[1], valBuf[2], valBuf[3] = val, val, val, val
	needle := archsimd.LoadInt32x4(&valBuf)

	for ; i+int32Lane <= n; i += int32Lane {
		v := archsimd.LoadInt32x4Slice(s[i:])
		mask := v.Equal(needle)
		b := mask.ToBits()
		if b != 0 {
			return i + bits.TrailingZeros8(b)
		}
	}

	for ; i < n; i++ {
		if s[i] == val {
			return i
		}
	}
	return -1
}

func findInt64(s []int64, val int64) int {
	n := len(s)
	i := 0

	var valBuf [2]int64
	valBuf[0], valBuf[1] = val, val
	needle := archsimd.LoadInt64x2(&valBuf)

	for ; i+int64Lane <= n; i += int64Lane {
		v := archsimd.LoadInt64x2Slice(s[i:])
		mask := v.Equal(needle)
		b := mask.ToBits()
		if b != 0 {
			return i + bits.TrailingZeros8(b)
		}
	}

	for ; i < n; i++ {
		if s[i] == val {
			return i
		}
	}
	return -1
}

func findFloat32(s []float32, val float32) int {
	n := len(s)
	i := 0

	var valBuf [4]float32
	valBuf[0], valBuf[1], valBuf[2], valBuf[3] = val, val, val, val
	needle := archsimd.LoadFloat32x4(&valBuf)

	for ; i+f32Lane <= n; i += f32Lane {
		v := archsimd.LoadFloat32x4Slice(s[i:])
		mask := v.Equal(needle)
		b := mask.ToBits()
		if b != 0 {
			return i + bits.TrailingZeros8(b)
		}
	}

	for ; i < n; i++ {
		if s[i] == val {
			return i
		}
	}
	return -1
}

func findFloat64(s []float64, val float64) int {
	n := len(s)
	i := 0

	var valBuf [2]float64
	valBuf[0], valBuf[1] = val, val
	needle := archsimd.LoadFloat64x2(&valBuf)

	for ; i+f64Lane <= n; i += f64Lane {
		v := archsimd.LoadFloat64x2Slice(s[i:])
		mask := v.Equal(needle)
		b := mask.ToBits()
		if b != 0 {
			return i + bits.TrailingZeros8(b)
		}
	}

	for ; i < n; i++ {
		if s[i] == val {
			return i
		}
	}
	return -1
}

// --- Count ---

func countInt32(s []int32, val int32) int {
	n := len(s)
	i := 0
	count := 0

	var valBuf [4]int32
	valBuf[0], valBuf[1], valBuf[2], valBuf[3] = val, val, val, val
	needle := archsimd.LoadInt32x4(&valBuf)

	for ; i+int32Lane <= n; i += int32Lane {
		v := archsimd.LoadInt32x4Slice(s[i:])
		mask := v.Equal(needle)
		count += bits.OnesCount8(mask.ToBits())
	}

	for ; i < n; i++ {
		if s[i] == val {
			count++
		}
	}
	return count
}

func countInt64(s []int64, val int64) int {
	n := len(s)
	i := 0
	count := 0

	var valBuf [2]int64
	valBuf[0], valBuf[1] = val, val
	needle := archsimd.LoadInt64x2(&valBuf)

	for ; i+int64Lane <= n; i += int64Lane {
		v := archsimd.LoadInt64x2Slice(s[i:])
		mask := v.Equal(needle)
		count += bits.OnesCount8(mask.ToBits())
	}

	for ; i < n; i++ {
		if s[i] == val {
			count++
		}
	}
	return count
}

func countFloat32(s []float32, val float32) int {
	n := len(s)
	i := 0
	count := 0

	var valBuf [4]float32
	valBuf[0], valBuf[1], valBuf[2], valBuf[3] = val, val, val, val
	needle := archsimd.LoadFloat32x4(&valBuf)

	for ; i+f32Lane <= n; i += f32Lane {
		v := archsimd.LoadFloat32x4Slice(s[i:])
		mask := v.Equal(needle)
		count += bits.OnesCount8(mask.ToBits())
	}

	for ; i < n; i++ {
		if s[i] == val {
			count++
		}
	}
	return count
}

func countFloat64(s []float64, val float64) int {
	n := len(s)
	i := 0
	count := 0

	var valBuf [2]float64
	valBuf[0], valBuf[1] = val, val
	needle := archsimd.LoadFloat64x2(&valBuf)

	for ; i+f64Lane <= n; i += f64Lane {
		v := archsimd.LoadFloat64x2Slice(s[i:])
		mask := v.Equal(needle)
		count += bits.OnesCount8(mask.ToBits())
	}

	for ; i < n; i++ {
		if s[i] == val {
			count++
		}
	}
	return count
}

// --- Min ---

func minInt32(s []int32) int32 {
	n := len(s)
	if n < int32Lane {
		m := s[0]
		for _, v := range s[1:] {
			if v < m {
				m = v
			}
		}
		return m
	}

	acc := archsimd.LoadInt32x4Slice(s)
	i := int32Lane
	for ; i+int32Lane <= n; i += int32Lane {
		v := archsimd.LoadInt32x4Slice(s[i:])
		acc = acc.Min(v)
	}

	var buf [4]int32
	acc.Store(&buf)
	m := buf[0]
	for _, v := range buf[1:] {
		if v < m {
			m = v
		}
	}
	for ; i < n; i++ {
		if s[i] < m {
			m = s[i]
		}
	}
	return m
}

func minInt64(s []int64) int64 {
	n := len(s)
	if n < int64Lane {
		return s[0]
	}

	acc := archsimd.LoadInt64x2Slice(s)
	i := int64Lane
	for ; i+int64Lane <= n; i += int64Lane {
		v := archsimd.LoadInt64x2Slice(s[i:])
		acc = acc.Min(v)
	}

	var buf [2]int64
	acc.Store(&buf)
	m := buf[0]
	if buf[1] < m {
		m = buf[1]
	}
	for ; i < n; i++ {
		if s[i] < m {
			m = s[i]
		}
	}
	return m
}

func minFloat32(s []float32) float32 {
	n := len(s)
	if n < f32Lane {
		m := s[0]
		for _, v := range s[1:] {
			if v < m {
				m = v
			}
		}
		return m
	}

	acc := archsimd.LoadFloat32x4Slice(s)
	i := f32Lane
	for ; i+f32Lane <= n; i += f32Lane {
		v := archsimd.LoadFloat32x4Slice(s[i:])
		acc = acc.Min(v)
	}

	var buf [4]float32
	acc.Store(&buf)
	m := buf[0]
	for _, v := range buf[1:] {
		if v < m {
			m = v
		}
	}
	for ; i < n; i++ {
		if s[i] < m {
			m = s[i]
		}
	}
	return m
}

func minFloat64(s []float64) float64 {
	n := len(s)
	if n < f64Lane {
		return s[0]
	}

	acc := archsimd.LoadFloat64x2Slice(s)
	i := f64Lane
	for ; i+f64Lane <= n; i += f64Lane {
		v := archsimd.LoadFloat64x2Slice(s[i:])
		acc = acc.Min(v)
	}

	var buf [2]float64
	acc.Store(&buf)
	m := buf[0]
	if buf[1] < m {
		m = buf[1]
	}
	for ; i < n; i++ {
		if s[i] < m {
			m = s[i]
		}
	}
	return m
}

// --- Max ---

func maxInt32(s []int32) int32 {
	n := len(s)
	if n < int32Lane {
		m := s[0]
		for _, v := range s[1:] {
			if v > m {
				m = v
			}
		}
		return m
	}

	acc := archsimd.LoadInt32x4Slice(s)
	i := int32Lane
	for ; i+int32Lane <= n; i += int32Lane {
		v := archsimd.LoadInt32x4Slice(s[i:])
		acc = acc.Max(v)
	}

	var buf [4]int32
	acc.Store(&buf)
	m := buf[0]
	for _, v := range buf[1:] {
		if v > m {
			m = v
		}
	}
	for ; i < n; i++ {
		if s[i] > m {
			m = s[i]
		}
	}
	return m
}

func maxInt64(s []int64) int64 {
	n := len(s)
	if n < int64Lane {
		return s[0]
	}

	acc := archsimd.LoadInt64x2Slice(s)
	i := int64Lane
	for ; i+int64Lane <= n; i += int64Lane {
		v := archsimd.LoadInt64x2Slice(s[i:])
		acc = acc.Max(v)
	}

	var buf [2]int64
	acc.Store(&buf)
	m := buf[0]
	if buf[1] > m {
		m = buf[1]
	}
	for ; i < n; i++ {
		if s[i] > m {
			m = s[i]
		}
	}
	return m
}

func maxFloat32(s []float32) float32 {
	n := len(s)
	if n < f32Lane {
		m := s[0]
		for _, v := range s[1:] {
			if v > m {
				m = v
			}
		}
		return m
	}

	acc := archsimd.LoadFloat32x4Slice(s)
	i := f32Lane
	for ; i+f32Lane <= n; i += f32Lane {
		v := archsimd.LoadFloat32x4Slice(s[i:])
		acc = acc.Max(v)
	}

	var buf [4]float32
	acc.Store(&buf)
	m := buf[0]
	for _, v := range buf[1:] {
		if v > m {
			m = v
		}
	}
	for ; i < n; i++ {
		if s[i] > m {
			m = s[i]
		}
	}
	return m
}

func maxFloat64(s []float64) float64 {
	n := len(s)
	if n < f64Lane {
		return s[0]
	}

	acc := archsimd.LoadFloat64x2Slice(s)
	i := f64Lane
	for ; i+f64Lane <= n; i += f64Lane {
		v := archsimd.LoadFloat64x2Slice(s[i:])
		acc = acc.Max(v)
	}

	var buf [2]float64
	acc.Store(&buf)
	m := buf[0]
	if buf[1] > m {
		m = buf[1]
	}
	for ; i < n; i++ {
		if s[i] > m {
			m = s[i]
		}
	}
	return m
}

// --- MinMax ---

func minmaxInt32(s []int32) (int32, int32) {
	n := len(s)
	if n < int32Lane {
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

	loAcc := archsimd.LoadInt32x4Slice(s)
	hiAcc := loAcc
	i := int32Lane
	for ; i+int32Lane <= n; i += int32Lane {
		v := archsimd.LoadInt32x4Slice(s[i:])
		loAcc = loAcc.Min(v)
		hiAcc = hiAcc.Max(v)
	}

	var loBuf, hiBuf [4]int32
	loAcc.Store(&loBuf)
	hiAcc.Store(&hiBuf)

	lo := loBuf[0]
	hi := hiBuf[0]
	for j := 1; j < 4; j++ {
		if loBuf[j] < lo {
			lo = loBuf[j]
		}
		if hiBuf[j] > hi {
			hi = hiBuf[j]
		}
	}
	for ; i < n; i++ {
		if s[i] < lo {
			lo = s[i]
		}
		if s[i] > hi {
			hi = s[i]
		}
	}
	return lo, hi
}

func minmaxInt64(s []int64) (int64, int64) {
	n := len(s)
	if n < int64Lane {
		return s[0], s[0]
	}

	loAcc := archsimd.LoadInt64x2Slice(s)
	hiAcc := loAcc
	i := int64Lane
	for ; i+int64Lane <= n; i += int64Lane {
		v := archsimd.LoadInt64x2Slice(s[i:])
		loAcc = loAcc.Min(v)
		hiAcc = hiAcc.Max(v)
	}

	var loBuf, hiBuf [2]int64
	loAcc.Store(&loBuf)
	hiAcc.Store(&hiBuf)

	lo, hi := loBuf[0], hiBuf[0]
	if loBuf[1] < lo {
		lo = loBuf[1]
	}
	if hiBuf[1] > hi {
		hi = hiBuf[1]
	}
	for ; i < n; i++ {
		if s[i] < lo {
			lo = s[i]
		}
		if s[i] > hi {
			hi = s[i]
		}
	}
	return lo, hi
}

func minmaxFloat32(s []float32) (float32, float32) {
	n := len(s)
	if n < f32Lane {
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

	loAcc := archsimd.LoadFloat32x4Slice(s)
	hiAcc := loAcc
	i := f32Lane
	for ; i+f32Lane <= n; i += f32Lane {
		v := archsimd.LoadFloat32x4Slice(s[i:])
		loAcc = loAcc.Min(v)
		hiAcc = hiAcc.Max(v)
	}

	var loBuf, hiBuf [4]float32
	loAcc.Store(&loBuf)
	hiAcc.Store(&hiBuf)

	lo, hi := loBuf[0], hiBuf[0]
	for j := 1; j < 4; j++ {
		if loBuf[j] < lo {
			lo = loBuf[j]
		}
		if hiBuf[j] > hi {
			hi = hiBuf[j]
		}
	}
	for ; i < n; i++ {
		if s[i] < lo {
			lo = s[i]
		}
		if s[i] > hi {
			hi = s[i]
		}
	}
	return lo, hi
}

func minmaxFloat64(s []float64) (float64, float64) {
	n := len(s)
	if n < f64Lane {
		return s[0], s[0]
	}

	loAcc := archsimd.LoadFloat64x2Slice(s)
	hiAcc := loAcc
	i := f64Lane
	for ; i+f64Lane <= n; i += f64Lane {
		v := archsimd.LoadFloat64x2Slice(s[i:])
		loAcc = loAcc.Min(v)
		hiAcc = hiAcc.Max(v)
	}

	var loBuf, hiBuf [2]float64
	loAcc.Store(&loBuf)
	hiAcc.Store(&hiBuf)

	lo, hi := loBuf[0], hiBuf[0]
	if loBuf[1] < lo {
		lo = loBuf[1]
	}
	if hiBuf[1] > hi {
		hi = hiBuf[1]
	}
	for ; i < n; i++ {
		if s[i] < lo {
			lo = s[i]
		}
		if s[i] > hi {
			hi = s[i]
		}
	}
	return lo, hi
}

// --- Element-wise Arithmetic ---

func addSlicesInt32(dst, s1, s2 []int32) {
	n := len(dst)
	i := 0
	for ; i+int32Lane <= n; i += int32Lane {
		v1 := archsimd.LoadInt32x4Slice(s1[i:])
		v2 := archsimd.LoadInt32x4Slice(s2[i:])
		v1.Add(v2).StoreSlice(dst[i:])
	}
	for ; i < n; i++ {
		dst[i] = s1[i] + s2[i]
	}
}

func addSlicesInt64(dst, s1, s2 []int64) {
	n := len(dst)
	i := 0
	for ; i+int64Lane <= n; i += int64Lane {
		v1 := archsimd.LoadInt64x2Slice(s1[i:])
		v2 := archsimd.LoadInt64x2Slice(s2[i:])
		v1.Add(v2).StoreSlice(dst[i:])
	}
	for ; i < n; i++ {
		dst[i] = s1[i] + s2[i]
	}
}

func addSlicesFloat32(dst, s1, s2 []float32) {
	n := len(dst)
	i := 0
	for ; i+f32Lane <= n; i += f32Lane {
		v1 := archsimd.LoadFloat32x4Slice(s1[i:])
		v2 := archsimd.LoadFloat32x4Slice(s2[i:])
		v1.Add(v2).StoreSlice(dst[i:])
	}
	for ; i < n; i++ {
		dst[i] = s1[i] + s2[i]
	}
}

func addSlicesFloat64(dst, s1, s2 []float64) {
	n := len(dst)
	i := 0
	for ; i+f64Lane <= n; i += f64Lane {
		v1 := archsimd.LoadFloat64x2Slice(s1[i:])
		v2 := archsimd.LoadFloat64x2Slice(s2[i:])
		v1.Add(v2).StoreSlice(dst[i:])
	}
	for ; i < n; i++ {
		dst[i] = s1[i] + s2[i]
	}
}

// mulSlicesInt64 is scalar: no SSE/AVX2 int64 multiply instruction.
// Exposed here so the generic dispatch can reach it without boxing.
func mulSlicesInt64(dst, s1, s2 []int64) {
	for i := range dst {
		dst[i] = s1[i] * s2[i]
	}
}

func mulSlicesInt32(dst, s1, s2 []int32) {
	n := len(dst)
	i := 0
	for ; i+int32Lane <= n; i += int32Lane {
		v1 := archsimd.LoadInt32x4Slice(s1[i:])
		v2 := archsimd.LoadInt32x4Slice(s2[i:])
		v1.Mul(v2).StoreSlice(dst[i:])
	}
	for ; i < n; i++ {
		dst[i] = s1[i] * s2[i]
	}
}

func mulSlicesFloat32(dst, s1, s2 []float32) {
	n := len(dst)
	i := 0
	for ; i+f32Lane <= n; i += f32Lane {
		v1 := archsimd.LoadFloat32x4Slice(s1[i:])
		v2 := archsimd.LoadFloat32x4Slice(s2[i:])
		v1.Mul(v2).StoreSlice(dst[i:])
	}
	for ; i < n; i++ {
		dst[i] = s1[i] * s2[i]
	}
}

func mulSlicesFloat64(dst, s1, s2 []float64) {
	n := len(dst)
	i := 0
	for ; i+f64Lane <= n; i += f64Lane {
		v1 := archsimd.LoadFloat64x2Slice(s1[i:])
		v2 := archsimd.LoadFloat64x2Slice(s2[i:])
		v1.Mul(v2).StoreSlice(dst[i:])
	}
	for ; i < n; i++ {
		dst[i] = s1[i] * s2[i]
	}
}

// --- Dot Product ---

func dotProductInt32(s1, s2 []int32) int32 {
	n := len(s1)
	i := 0
	var acc archsimd.Int32x4

	for ; i+int32Lane <= n; i += int32Lane {
		v1 := archsimd.LoadInt32x4Slice(s1[i:])
		v2 := archsimd.LoadInt32x4Slice(s2[i:])
		acc = acc.Add(v1.Mul(v2))
	}

	var buf [4]int32
	acc.Store(&buf)
	total := buf[0] + buf[1] + buf[2] + buf[3]

	for ; i < n; i++ {
		total += s1[i] * s2[i]
	}
	return total
}

// dotProductInt64 is intentionally scalar: SSE/AVX2 have no 64-bit integer
// multiply (PMULLQ exists only on AVX-512 DQ). When the project gains an
// AVX-512 dispatch this can become vectorized; until then it matches the
// hand-written loop exactly so callers see no surprise.
func dotProductInt64(s1, s2 []int64) int64 {
	var total int64
	for i := range s1 {
		total += s1[i] * s2[i]
	}
	return total
}

func dotProductFloat32(s1, s2 []float32) float32 {
	n := len(s1)
	i := 0
	var acc archsimd.Float32x4

	for ; i+f32Lane <= n; i += f32Lane {
		v1 := archsimd.LoadFloat32x4Slice(s1[i:])
		v2 := archsimd.LoadFloat32x4Slice(s2[i:])
		acc = acc.Add(v1.Mul(v2))
	}

	var buf [4]float32
	acc.Store(&buf)
	total := buf[0] + buf[1] + buf[2] + buf[3]

	for ; i < n; i++ {
		total += s1[i] * s2[i]
	}
	return total
}

func dotProductFloat64(s1, s2 []float64) float64 {
	n := len(s1)
	i := 0
	var acc archsimd.Float64x2

	for ; i+f64Lane <= n; i += f64Lane {
		v1 := archsimd.LoadFloat64x2Slice(s1[i:])
		v2 := archsimd.LoadFloat64x2Slice(s2[i:])
		acc = acc.Add(v1.Mul(v2))
	}

	var buf [2]float64
	acc.Store(&buf)
	total := buf[0] + buf[1]

	for ; i < n; i++ {
		total += s1[i] * s2[i]
	}
	return total
}

