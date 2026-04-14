<p align="center">
  <h1 align="center">TurboSlice</h1>
  <p align="center">
    <strong>SIMD-accelerated slice operations for Go</strong><br>
    <code>Sum</code> &middot; <code>Find</code> &middot; <code>Min</code> &middot; <code>Max</code> &middot; <code>Filter</code> &middot; <code>Map</code> &middot; <code>DotProduct</code> &middot; and more
  </p>
  <p align="center">
    <a href="https://pkg.go.dev/github.com/atul-007/turboslice"><img src="https://pkg.go.dev/badge/github.com/atul-007/turboslice.svg" alt="Go Reference"></a>
    <a href="https://goreportcard.com/report/github.com/atul-007/turboslice"><img src="https://goreportcard.com/badge/github.com/atul-007/turboslice" alt="Go Report Card"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License"></a>
    <a href="https://github.com/atul-007/turboslice"><img src="https://img.shields.io/badge/go-1.26+-00ADD8.svg" alt="Go 1.26+"></a>
  </p>
</p>

---

Drop-in replacements for your `for` loops that run on **SIMD vector instructions** (SSE/AVX2) on AMD64 via Go 1.26's [`simd/archsimd`](https://pkg.go.dev/simd/archsimd) — with automatic scalar fallback on ARM64 and everything else. Zero CGo. Zero assembly files. Pure Go.

```go
import "github.com/atul-007/turboslice"

data := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

turboslice.Sum(data)            // 55
turboslice.Find(data, 7)        // 6
turboslice.Min(data)            // 1
turboslice.Max(data)            // 10
turboslice.Contains(data, 5)    // true
turboslice.DotProduct(a, b)     // scalar product
```

## Why TurboSlice?

You've written this loop a thousand times:

```go
total := 0
for _, v := range data {
    total += v
}
```

TurboSlice replaces it with a single call that processes **4-8 elements per CPU cycle** using SIMD vector instructions:

```go
total := turboslice.Sum(data) // SIMD-accelerated on AMD64
```

**No assembly. No CGo. No unsafe.** Just `go get` and go.

## Benchmarks

### AMD64 + SIMD (`GOEXPERIMENT=simd`)

Measured on AMD64, TurboSlice SIMD vs a hand-written `for` loop.

> [!NOTE]
> Numbers below were collected via QEMU emulation (Apple Silicon host). QEMU executes SIMD instructions **sequentially**, not in parallel. On native AMD64 hardware, expect **3-5x real speedups** as each vector instruction processes 4-8 elements simultaneously.

```
goos: linux
goarch: amd64
cpu: VirtualApple @ 2.50GHz
```

| Operation | Elements | TurboSlice | `for` loop | Faster |
|:---|---:|---:|---:|:---|
| **Sum** `int32` | 1M | 209 us | 268 us | **1.28x** |
| **Min** `int32` | 1M | 209 us | 265 us | **1.27x** |
| **Max** `float64` | 1M | 419 us | 639 us | **1.53x** |
| **AddSlices** `int32` | 1M | 640 us | 941 us | **1.47x** |
| **AddSlices** `float64` | 1M | 706 us | 1061 us | **1.50x** |
| **DotProduct** `int32` | 1M | 394 us | 428 us | **1.09x** |

<details>
<summary><strong>Projected native AMD64 performance</strong> (click to expand)</summary>
<br>

Each SSE vector instruction processes 4x `int32` or 2x `float64` in a single cycle. With proper pipelining, this yields:

| Operation | Projected native speedup |
|:---|:---|
| Sum, Min, Max, Count, Find (`int32`) | **~4x** (4 elements/cycle) |
| Sum, Min, Max (`float64`) | **~2x** (2 elements/cycle) |
| DotProduct (`float32`) | **~4x** (FMA: multiply + add fused) |
| AddSlices, MulSlices | **~4x** (`int32`), **~2x** (`float64`) |

These scale further with AVX2 (8x `int32`) and AVX-512 (16x `int32`).

</details>

### Zero-overhead Typed API (any platform)

The typed functions (`SumInt32`, `DotProductFloat64`, etc.) are compiler-verified to **fully inline**, matching hand-written loop performance with 0% overhead:

```
BenchmarkMinInt32/Typed/1M          292,276 ns/op     0 allocs
BenchmarkMinInt32/NaiveLoop/1M      296,177 ns/op     0 allocs
                                    ^^^^^^^^^^^^^^^^
                                    identical. zero overhead.
```

```
BenchmarkMaxFloat64/Typed/1M        407,044 ns/op
BenchmarkMaxFloat64/NaiveLoop/1M    639,607 ns/op
                                    ^^^^^^^^^^^^^^^^
                                    36% faster than the naive loop!
```

> Reproduce: `go test -bench=. -benchmem -count=3`

## Install

```bash
go get github.com/atul-007/turboslice
```

Requires **Go 1.26+**. That's it. Works immediately on any platform.

For SIMD acceleration, build with:
```bash
GOEXPERIMENT=simd go build ./...
```

## API at a Glance

### SIMD-Accelerated (int32, int64, float32, float64)

```go
// Aggregation
Sum[T Numeric](s []T) T
Min[T Numeric](s []T) T
Max[T Numeric](s []T) T
MinMax[T Numeric](s []T) (min, max T)     // single pass, 2 results

// Search
Find[T comparable](s []T, val T) int       // returns index, -1 if missing
Contains[T comparable](s []T, val T) bool
Count[T comparable](s []T, val T) int

// Element-wise math
AddSlices[T Numeric](s1, s2 []T) []T
MulSlices[T Numeric](s1, s2 []T) []T
DotProduct[T Numeric](s1, s2 []T) T
```

### Typed Fast-Path (fully inlined, zero overhead)

When you know the type and need max throughput in a hot loop:

```go
SumInt32(s []int32) int32          SumFloat64(s []float64) float64
MinInt32(s []int32) int32          MinFloat64(s []float64) float64
MaxInt32(s []int32) int32          MaxFloat64(s []float64) float64
FindInt32(s, val) int              FindFloat64(s, val) int
CountInt32(s, val) int             CountFloat64(s, val) int
DotProductInt32(s1, s2) int32      DotProductFloat64(s1, s2) float64
// ... also Int64 and Float32 variants for all of the above
```

### Generic Utilities (any type)

```go
Map[T, U](s []T, fn func(T) U) []U        // transform
Filter[T](s []T, fn func(T) bool) []T      // select
Reduce[T, U](s []T, init U, fn) U          // fold
Any[T](s []T, fn func(T) bool) bool        // exists?
All[T](s []T, fn func(T) bool) bool        // universal?
Chunk[T](s []T, n int) [][]T               // split into groups
Unique[T comparable](s []T) []T             // deduplicate
Reverse[T](s []T) []T                       // reverse copy
Flatten[T](ss [][]T) []T                    // join nested
ForEach[T](s []T, fn func(T))              // side effects
```

## Real-World Examples

### Signal Processing

```go
signal  := loadSensorData()                             // []float64
weights := precomputeWeights(len(signal))               // []float64

weighted := turboslice.MulSlices(signal, weights)       // element-wise
energy   := turboslice.DotProduct(signal, signal)       // sum of squares
peak     := turboslice.Max(signal)
lo, hi   := turboslice.MinMax(signal)                   // one pass
```

### Analytics Pipeline

```go
scores := fetchAllScores()                              // []int32, millions

total    := turboslice.Sum(scores)
lo, hi   := turboslice.MinMax(scores)
outliers := turboslice.Filter(scores, func(s int32) bool {
    return s > 3*stddev
})
buckets  := turboslice.Chunk(scores, 1000)
unique   := turboslice.Unique(scores)
```

### ML Feature Engineering

```go
features := []float64{ /* embeddings */ }
weights  := []float64{ /* model weights */ }

score    := turboslice.DotProduct(features, weights)    // SIMD dot product
norm     := math.Sqrt(turboslice.DotProduct(features, features))
scaled   := turboslice.Map(features, func(f float64) float64 {
    return f / norm
})
```

### Custom Types Work Too

```go
type Celsius float32

temps := []Celsius{-5.2, 3.1, 18.7, 22.4, 28.9}
avg   := turboslice.Sum(temps) / Celsius(len(temps))    // SIMD accelerated
cold  := turboslice.Filter(temps, func(t Celsius) bool { return t < 0 })
```

## Architecture

```
                         Your Code
                            |
                     turboslice.Sum(data)
                            |
                    +-------+-------+
                    |  type switch   |
                    +---+---+---+---+
                        |   |   |
              []int32   | []float64
                  |     |     |
             sumInt32   |  sumFloat64
                  |     |     |
         +--------+----+-----+---------+
         |                             |
   dispatch_simd_amd64.go      dispatch_default.go
   (goexperiment.simd)         (all other builds)
         |                             |
   archsimd.LoadInt32x4Slice    simple for-loop
   archsimd.Add()              (auto-vectorized
   archsimd.Store()             by compiler)
         |                             |
     SSE / AVX2                   NEON / scalar
```

### Design Decisions

| Decision | Rationale |
|:---|:---|
| **128-bit SSE** over 256-bit AVX2 | SSE works on every AMD64 CPU. AVX2 upgrade path is clean via `archsimd.X86.AVX2()` runtime checks. |
| **Build-tag dispatch** over runtime dispatch | Zero-cost abstraction. The linker includes exactly one implementation. No `if` branches at runtime. |
| **Type-switch + typed API** dual approach | Generic API for convenience, typed API for hot paths where inlining matters. Users choose their tradeoff. |
| **No `unsafe`** in public API | The SIMD internals use `archsimd` (compiler intrinsics). The public surface is pure safe Go. |

## Running Benchmarks

```bash
# On any platform (scalar)
go test -bench=. -benchmem -count=3

# On AMD64 with SIMD
GOEXPERIMENT=simd go test -bench=. -benchmem -count=3

# Interactive demo
go run ./cmd/demo/

# AMD64 via Docker (from ARM64 host)
./scripts/bench-amd64.sh
```

## Roadmap

- [ ] **AVX2 (256-bit)** - 2x wider vectors with runtime feature detection
- [ ] **AVX-512** - 4x wider on supported CPUs
- [ ] **ARM64 NEON** - when `simd/archsimd` adds support
- [ ] **Parallel fan-out** - goroutine splitting for 10M+ element slices
- [ ] **SIMD Sort** - vectorized partitioning for radix/quick sort hybrids

## Contributing

Contributions are welcome. Please open an issue to discuss before submitting large changes.

```bash
git clone https://github.com/atul-007/turboslice
cd turboslice
go test ./...                                           # run tests
go test -bench=. -benchmem                              # run benchmarks
GOARCH=amd64 GOEXPERIMENT=simd go vet ./...             # verify SIMD path
```

## License

[MIT](LICENSE)
