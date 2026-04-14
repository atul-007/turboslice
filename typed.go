package turboslice

// Direct-typed functions that bypass generic dispatch overhead.
// Use these in hot paths where every nanosecond matters.
// On ARM64, these are ~2-5x faster than the generic versions
// because the compiler can fully inline and auto-vectorize the loops.
// On AMD64+SIMD, both generic and typed versions use SIMD acceleration.

// --- Sum ---

func SumInt32(s []int32) int32     { return sumInt32(s) }
func SumInt64(s []int64) int64     { return sumInt64(s) }
func SumFloat32(s []float32) float32 { return sumFloat32(s) }
func SumFloat64(s []float64) float64 { return sumFloat64(s) }

// --- Find ---

func FindInt32(s []int32, val int32) int       { return findInt32(s, val) }
func FindInt64(s []int64, val int64) int       { return findInt64(s, val) }
func FindFloat32(s []float32, val float32) int { return findFloat32(s, val) }
func FindFloat64(s []float64, val float64) int { return findFloat64(s, val) }

// --- Count ---

func CountInt32(s []int32, val int32) int       { return countInt32(s, val) }
func CountInt64(s []int64, val int64) int       { return countInt64(s, val) }
func CountFloat32(s []float32, val float32) int { return countFloat32(s, val) }
func CountFloat64(s []float64, val float64) int { return countFloat64(s, val) }

// --- Min ---

func MinInt32(s []int32) int32     { return minInt32(s) }
func MinInt64(s []int64) int64     { return minInt64(s) }
func MinFloat32(s []float32) float32 { return minFloat32(s) }
func MinFloat64(s []float64) float64 { return minFloat64(s) }

// --- Max ---

func MaxInt32(s []int32) int32     { return maxInt32(s) }
func MaxInt64(s []int64) int64     { return maxInt64(s) }
func MaxFloat32(s []float32) float32 { return maxFloat32(s) }
func MaxFloat64(s []float64) float64 { return maxFloat64(s) }

// --- DotProduct ---

func DotProductInt32(s1, s2 []int32) int32       { return dotProductInt32(s1, s2) }
func DotProductInt64(s1, s2 []int64) int64       { return dotProductInt64(s1, s2) }
func DotProductFloat32(s1, s2 []float32) float32 { return dotProductFloat32(s1, s2) }
func DotProductFloat64(s1, s2 []float64) float64 { return dotProductFloat64(s1, s2) }
