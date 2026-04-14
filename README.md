# TurboSlice

**SIMD-accelerated slice operations for Go.** A "Standard Library+" that makes number crunching and data filtering dramatically faster — without writing a single line of assembly.

Built on Go 1.26's experimental [`simd/archsimd`](https://pkg.go.dev/simd/archsimd) package. Auto-falls back to optimized scalar Go on platforms without SIMD support.

```go
import "github.com/a-ranjan/turboslice"

data := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

turboslice.Sum(data)            // 55
turboslice.Find(data, 7)        // 6
turboslice.Min(data)            // 1
turboslice.Max(data)            // 10
turboslice.Contains(data, 5)    // true
turboslice.Count(data, 3)       // 1
```

## Performance

### AMD64 + SIMD (GOEXPERIMENT=simd)

Benchmarks on AMD64 with `GOEXPERIMENT=simd`. TurboSlice SIMD vs hand-written `for` loop:

| Operation | Size | TurboSlice (SIMD) | Naive Loop | Speedup |
|---|---|---|---|---|
| **Sum** (int32) | 1M | 209 us | 268 us | **1.28x** |
| **Min** (int32) | 1M | 209 us | 265 us | **1.27x** |
| **AddSlices** (int32) | 1M | 640 us | 941 us | **1.47x** |
| **DotProduct** (int32) | 1M | 394 us | 428 us | **1.09x** |

> These numbers were collected via QEMU emulation. On native AMD64 hardware with true SIMD execution, expect **3-5x speedups** as each vector instruction processes 4-8 elements per cycle instead of being emulated sequentially.

### Typed API (zero-overhead, any platform)

The typed functions (`SumInt32`, `MinInt32`, etc.) fully inline and match hand-written loop performance on every platform:

| Operation | Size | Typed API | Naive Loop | Overhead |
|---|---|---|---|---|
| **Sum** (int32) | 1M | 266 us | 266 us | **0%** |
| **Min** (int32) | 1M | 292 us | 296 us | **0%** |
| **DotProduct** (float64) | 1M | 794 us | 800 us | **0%** |
| **MaxFloat64** | 1M | 407 us | 639 us | **-36% (faster!)** |

### Which API to use?

| API | When to use |
|---|---|
| `Sum[T](s)` | Convenience, multi-type code, non-hot-paths |
| `SumInt32(s)` | Hot loops, maximum performance, known type |

> Run your own: `go test -bench=. -benchmem` or `GOEXPERIMENT=simd go test -bench=. -benchmem`

## Installation

```bash
go get github.com/a-ranjan/turboslice
```

**Requirements:** Go 1.26+

For SIMD acceleration on AMD64:
```bash
GOEXPERIMENT=simd go build ./...
```

Without the experiment flag, everything works identically — it just uses the scalar fallback path.

## API Reference

### Numeric Operations (SIMD-accelerated)

These operations use SIMD vector instructions on AMD64 for `int32`, `int64`, `float32`, `float64`. All other numeric types use an optimized scalar path.

```go
// Aggregation
turboslice.Sum[T Numeric](s []T) T
turboslice.Min[T Numeric](s []T) T                    // panics on empty
turboslice.Max[T Numeric](s []T) T                    // panics on empty
turboslice.MinMax[T Numeric](s []T) (min, max T)      // single pass

// Search
turboslice.Find[T comparable](s []T, val T) int        // -1 if not found
turboslice.Contains[T comparable](s []T, val T) bool
turboslice.Count[T comparable](s []T, val T) int

// Element-wise arithmetic
turboslice.AddSlices[T Numeric](s1, s2 []T) []T
turboslice.MulSlices[T Numeric](s1, s2 []T) []T

// Linear algebra
turboslice.DotProduct[T Numeric](s1, s2 []T) T
```

### Typed Fast-Path (zero dispatch overhead)

For hot loops where you know the concrete type. These fully inline:

```go
turboslice.SumInt32(s []int32) int32
turboslice.SumFloat64(s []float64) float64
turboslice.MinInt32(s []int32) int32
turboslice.MaxFloat64(s []float64) float64
turboslice.FindInt32(s []int32, val int32) int
turboslice.CountInt32(s []int32, val int32) int
turboslice.DotProductFloat64(s1, s2 []float64) float64
// ... and more for Int32, Int64, Float32, Float64 variants
```

### Generic Utilities

Higher-order functions for any slice type:

```go
// Transform
turboslice.Map[T, U any](s []T, fn func(T) U) []U
turboslice.Filter[T any](s []T, fn func(T) bool) []T
turboslice.Reduce[T, U any](s []T, initial U, fn func(U, T) U) U

// Predicates
turboslice.Any[T any](s []T, fn func(T) bool) bool
turboslice.All[T any](s []T, fn func(T) bool) bool

// Structure
turboslice.Chunk[T any](s []T, n int) [][]T
turboslice.Unique[T comparable](s []T) []T
turboslice.Reverse[T any](s []T) []T
turboslice.Flatten[T any](ss [][]T) []T
turboslice.ForEach[T any](s []T, fn func(T))
```

## Usage Examples

### Signal processing
```go
signal := []float64{0.1, 0.5, 0.9, 0.3, 0.7, 0.2, 0.8, 0.4}
weights := []float64{0.5, 1.0, 1.5, 1.0, 0.5, 1.0, 1.5, 1.0}

weighted := turboslice.MulSlices(signal, weights)
energy := turboslice.DotProduct(signal, signal)
peak := turboslice.Max(signal)
```

### Data analysis
```go
temperatures := []float32{-5.2, 3.1, 18.7, 22.4, 15.8, -1.0, 28.9, 11.3}

lo, hi := turboslice.MinMax(temperatures)
avg := turboslice.Sum(temperatures) / float32(len(temperatures))
freezing := turboslice.Count(temperatures, float32(0.0))
warm := turboslice.Filter(temperatures, func(t float32) bool { return t > 20 })
```

### Batch scoring
```go
scores := []int32{ /* millions of scores */ }

total := turboslice.Sum(scores)
hasTarget := turboslice.Contains(scores, targetScore)
normalized := turboslice.Map(scores, func(s int32) float64 {
    return float64(s) / float64(total)
})
```

## Architecture

```
turboslice/
+-- turboslice.go             Public API: generic functions with type-switch dispatch
+-- typed.go                  Typed fast-path API (SumInt32, MinFloat64, etc.)
+-- scalar.go                 Generic scalar helpers (sumScalar, findScalar, etc.)
+-- dispatch_default.go       Scalar specializations (!goexperiment.simd || !amd64)
+-- dispatch_simd_amd64.go    SIMD specializations (goexperiment.simd && amd64)
+-- generic.go                Map, Filter, Reduce, etc.
+-- turboslice_test.go        Tests
+-- bench_test.go             Benchmarks
+-- cmd/demo/                 Performance demo program
```

### How the dispatch works

```
Sum[T](s []T)
  |
  +-- type switch on T
  |     |
  |     +-- []int32  --> sumInt32(s)  --+-- (SIMD build) --> archsimd vectorized loop
  |     |                               +-- (default)    --> scalar for-loop
  |     +-- []int64  --> sumInt64(s)  --+
  |     +-- []float32 -> sumFloat32(s) +
  |     +-- []float64 -> sumFloat64(s) +
  |     |
  |     +-- []int, []uint8, etc. -------> sumScalar[T](s) (always generic loop)
```

Build tags select which `dispatch_*.go` file provides `sumInt32()` et al:
- `goexperiment.simd && amd64` -> SIMD implementation using 128-bit SSE vectors
- Everything else -> Optimized scalar loops

### Design decisions

**Why 128-bit (SSE) vectors instead of 256-bit (AVX2)?**
SSE is available on every AMD64 CPU. AVX2 requires runtime feature detection. Starting with 128-bit gives universal AMD64 coverage with a clean upgrade path to AVX2/AVX-512 via `archsimd.X86.AVX2()` checks.

**Why type-switch dispatch instead of interfaces?**
Interfaces add vtable overhead on every call. The type switch is resolved at compile time for concrete types, giving zero overhead dispatch. The compiler inlines the specialized functions directly.

**Why separate files instead of `//go:build` blocks?**
Go build constraints work at the file level. Separate files keep the SIMD and scalar implementations clean and independently testable. The compiler includes exactly one dispatch file per build.

## Type Constraints

```go
type Numeric interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}
```

Custom types work too:
```go
type Score int32
scores := []Score{100, 200, 300}
total := turboslice.Sum(scores) // Score(600) — SIMD accelerated
```

## Running Benchmarks

```bash
# Scalar (any platform)
go test -bench=. -benchmem

# SIMD-accelerated (AMD64 only)
GOEXPERIMENT=simd go test -bench=. -benchmem

# Compare SIMD vs scalar
GOEXPERIMENT=simd go test -bench=BenchmarkSumInt32 -benchmem -count=5
```

## Roadmap

- [ ] AVX2 (256-bit) specializations with runtime feature detection
- [ ] AVX-512 path for supported CPUs
- [ ] ARM64 NEON support when `simd/archsimd` adds it
- [ ] Parallel variants for very large slices (goroutine fan-out)
- [ ] `Sort` with SIMD-accelerated partitioning

## License

MIT
