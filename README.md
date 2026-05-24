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

Drop-in replacements for your `for` loops that run on **128-bit SSE vector instructions** on AMD64 via Go 1.26's [`simd/archsimd`](https://pkg.go.dev/simd/archsimd) — with automatic scalar fallback on ARM64 and everything else. Zero CGo. Zero assembly files. Pure Go.

Measured on **Intel Xeon Platinum 8481C** (Sapphire Rapids), TurboSlice wins **24-62%** on aggregations over slices ≥64K elements. See [benchmarks](#benchmarks) and [when SIMD helps vs when it doesn't](#when-simd-helps-vs-when-it-doesnt) below for the honest scope.

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

TurboSlice replaces it with a single call that processes **4 `int32`s per SSE instruction** — and on Sapphire Rapids that turns into a **2.6x speedup** for `Min`/`Max` over 1M elements:

```go
total := turboslice.Sum(data) // SIMD-accelerated on AMD64 for slices >=4K
```

**No assembly. No CGo. No unsafe.** Just `go get` and go.

## Benchmarks

### AMD64 + SIMD (`GOEXPERIMENT=simd`)

Native run on **GCE c3-standard-4** — Intel Xeon Platinum 8481C @ 2.7 GHz
(Sapphire Rapids, AVX-512), Linux 6.1, Go 1.26.1, performance governor,
6 iterations via `benchstat`. Reproduce with `./scripts/bench-native.sh`
on any linux/amd64 box.

```
goos: linux
goarch: amd64
cpu: Intel(R) Xeon(R) Platinum 8481C CPU @ 2.70GHz
```

| Operation | Elements | TurboSlice (SIMD) | Scalar build | Speedup |
|:---|---:|---:|---:|:---|
| **Min** `int32` | 1M | 267 µs | 705 µs | **2.64x** |
| **Min** `int32` | 64K | 16.5 µs | 44.0 µs | **2.66x** |
| **Min** `int32` | 1K | 275 ns | 681 ns | **2.48x** |
| **Count** `int32` | 1M | 360 µs | 681 µs | **1.89x** |
| **Count** `int32` | 64K | 22.5 µs | 38.7 µs | **1.72x** |
| **Sum** `int32` | 1M | 266 µs | 353 µs | **1.33x** |
| **Sum** `int32` | 64K | 16.7 µs | 22.1 µs | **1.32x** |
| **Sum** `float64` | 1M | 531 µs | 705 µs | **1.33x** |
| **DotProduct** `int32` | 1M | 453 µs | 700 µs | **1.55x** |
| **AddSlices** `int32` | 1M | 701 µs | 878 µs | **1.25x** |

**Geomean across the whole benchmark suite: −6.30%** (SIMD build vs scalar build).

### When SIMD helps vs when it doesn't

Hand-written 128-bit SSE doesn't beat the Go compiler's auto-vectorizer
everywhere. On Sapphire Rapids the compiler already emits AVX-512 for trivial
reduction loops; the SIMD path only wins where lane-comparison patterns
(min/max/count) aren't autovectorized as well. The breakdown:

| Operation | What runs on the SIMD build |
|:---|:---|
| `Min`, `Max`, `MinMax`, `Count` | **SSE** — 2-3x win at every size |
| `Sum`, `DotProduct[int32/float32]` | **SSE for N ≥ 4K (Sum) / N ≥ 16K (DotProduct), scalar otherwise** — overhead doesn't amortize at small N |
| `AddSlices[int32/float32]`, `MulSlices` | **SSE** — modest win on large slices |
| `Find`, `Contains`, `DotProduct[float64]`, `AddSlices[float64]` | **Scalar** — compiler autovectorization beats hand-written SSE here, so we don't try |
| `MulSlices[int64]`, `DotProduct[int64]` | **Scalar everywhere** — no 64-bit integer multiply in SSE/AVX2 (`PMULLQ` is AVX-512 DQ only) |

This is honest: the library detects when SIMD wouldn't help and falls through
to the scalar implementation, so you don't pay an overhead tax for small slices
or operations the compiler already vectorizes well.

### Size thresholds

The SIMD-build implementations of `Sum` and `DotProductInt32`/`DotProductFloat32`
guard on slice length and run a scalar loop below the crossover point:

```go
const (
    simdSumMinN = 4096   // Sum* uses SIMD for N >= 4096
    simdDotMinN = 16384  // DotProduct[int32/float32] uses SIMD for N >= 16384
)
```

These were chosen from the per-N benchmark data — below the threshold,
function-call overhead and SIMD setup cost outweighs the throughput gain.

### Typed API has zero dispatch overhead

The typed entry points (`SumInt32`, `MinInt32`, …) inline into the kernel
without going through `interface{}` or a type switch. On the scalar build they
match the naive hand-written loop within noise:

```
BenchmarkMinInt32/Typed/1M          704.9 µs    0 allocs
BenchmarkMinInt32/NaiveLoop/1M      704.3 µs    0 allocs
```

On the SIMD build they collect the win:

```
BenchmarkMinInt32/Typed/1M          266.6 µs    (2.64x faster)
```

> Reproduce: `./scripts/bench-native.sh` on a real Linux box, or
> `go test -bench=. -benchmem -count=3` for a quick local check.
> Numbers will vary with CPU vendor, frequency, and Go version.

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

## Behavior notes

A few things to know before reaching for these in production:

- **`Min`, `Max`, `MinMax` panic on empty slices.** Matches `slices.Min`/`Max`
  from the standard library. Check `len(s) > 0` first.
- **`Sum` returns the zero value for empty slices.** No panic.
- **`AddSlices`, `MulSlices`, `DotProduct` silently truncate** to
  `min(len(s1), len(s2))`. Pass equal-length slices to avoid surprises.
- **Integer overflow is not checked.** `Sum` and `DotProduct` accumulate in
  the element type, same as a hand-written loop, so wide reductions on
  narrow integers (`int8`, `int16`, `int32`) can wrap. Cast to a wider type
  before summing if that's a concern.
- **`int64` multiplication is scalar even on SIMD builds.** SSE/AVX2 have no
  64-bit integer multiply (`PMULLQ` is AVX-512 DQ only). `MulSlices[int64]`
  and `DotProduct[int64]` therefore run a plain `for` loop on every
  architecture.
- **`Find`, `Contains`, `DotProduct[float64]`, `AddSlices[float64]` are
  scalar on every build.** The Go compiler's auto-vectorizer produces
  faster code than hand-written 128-bit SSE for these patterns on modern
  AMD64 (especially CPUs with AVX2/AVX-512), so the SIMD build deliberately
  doesn't override them.
- **`Sum` and `DotProduct[int32/float32]` use scalar below a size
  threshold.** SIMD only activates for `Sum` at N ≥ 4096 and for these
  `DotProduct` variants at N ≥ 16384. Below that, function-call and SIMD
  setup overhead outweighs throughput. The threshold is a compile-time
  constant in [`dispatch_simd_amd64.go`](dispatch_simd_amd64.go).
- **NaN handling differs between SIMD and scalar `Min`/`Max`.** SSE
  `MINPS`/`MAXPS` returns the second operand when either side is NaN; the
  scalar fallback uses `<`/`>`, which is always false against NaN. Filter
  NaNs before calling if you need deterministic behavior.

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

# Native AMD64 (best numbers — runs both builds + benchstat)
./scripts/bench-native.sh           # ~7 min, default profile
./scripts/bench-native.sh quick     # ~90s smoke test
./scripts/bench-native.sh deep      # ~20 min, publishable
```

### Native benchmark on a Linux box (for contributors)

If you have access to a real linux/amd64 machine, `scripts/bench-native.sh`
runs the full bench suite under both build paths, captures CPU/kernel/Go
fingerprint, and produces a single `bench-results/summary.md` ready to
attach to an issue or PR:

```bash
git clone https://github.com/atul-007/turboslice
cd turboslice
./scripts/bench-native.sh                          # collect
tar -czf bench-results.tgz bench-results/          # bundle
# attach bench-results.tgz to an issue
```

For low-variance numbers, set the CPU governor to `performance` first
(the script prints a warning if it isn't):

```bash
sudo cpupower frequency-set -g performance
```

## Roadmap

- [ ] **AVX2/AVX-512 intrinsics** — only for ops where the compiler's
  auto-vectorizer doesn't already use them (today that's `Min`/`Max`/`Count`
  on `int32`/`int64`; the float reductions are already AVX-512 in the
  scalar build). Gated on `archsimd.X86.AVX2()`/`AVX512F()`.
- [ ] **ARM64 NEON** — when `simd/archsimd` adds support.
- [ ] **Parallel fan-out** — goroutine splitting for 10M+ element slices.
- [ ] **SIMD Sort** — vectorized partitioning for radix/quick sort hybrids.

## Contributing

Contributions are welcome. Please open an issue to discuss before submitting large changes.

```bash
git clone https://github.com/atul-007/turboslice
cd turboslice
go test ./...                                           # run tests
go test -bench=. -benchmem                              # run benchmarks
go test -fuzz=FuzzSumInt32 -fuzztime=30s                # exercise tail loops
GOARCH=amd64 GOEXPERIMENT=simd go vet ./...             # verify SIMD path
```

## License

[MIT](LICENSE)
