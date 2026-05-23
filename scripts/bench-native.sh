#!/usr/bin/env bash
# Run TurboSlice benchmarks on native hardware and produce a single
# self-contained results bundle that can be emailed back.
#
# Best on linux/amd64 (where the SIMD path runs). On non-amd64 hosts the
# SIMD step is skipped automatically and only the scalar baseline runs.
#
# Usage:
#   ./scripts/bench-native.sh             # default profile (~7 min)
#   ./scripts/bench-native.sh quick       # ~90 s smoke test
#   ./scripts/bench-native.sh deep        # ~20 min, publishable
#
# Outputs to ./bench-results/ :
#   env.txt        — hostname, kernel, CPU info, Go version, governor
#   scalar.txt     — go test -bench output without SIMD
#   simd.txt       — go test -bench output with GOEXPERIMENT=simd
#   compare.txt    — benchstat scalar -> simd
#   summary.md     — single markdown file ready to paste in the README
#
# Email-back instruction:
#   tar -czf bench-results.tgz bench-results/
#   # send bench-results.tgz to atul

set -euo pipefail

cd "$(dirname "$0")/.."

mode="${1:-default}"
case "$mode" in
    quick)   COUNT=3;  TIME="300ms" ;;
    default) COUNT=6;  TIME="1s"    ;;
    deep)    COUNT=10; TIME="3s"    ;;
    *) echo "usage: $0 [quick|default|deep]" >&2; exit 1 ;;
esac

OUT_DIR="bench-results"
mkdir -p "$OUT_DIR"
SCALAR_OUT="$OUT_DIR/scalar.txt"
SIMD_OUT="$OUT_DIR/simd.txt"
ENV_OUT="$OUT_DIR/env.txt"
CMP_OUT="$OUT_DIR/compare.txt"
SUMMARY="$OUT_DIR/summary.md"

# --- preflight: go version ---

if ! command -v go >/dev/null 2>&1; then
    echo "ERROR: 'go' not found in PATH. Install Go 1.26+ from https://go.dev/dl/" >&2
    exit 1
fi

GO_VER=$(go version | awk '{print $3}' | sed 's/go//')
GO_MAJOR=$(echo "$GO_VER" | cut -d. -f1)
GO_MINOR=$(echo "$GO_VER" | cut -d. -f2)
if [ "$GO_MAJOR" -lt 1 ] || { [ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 26 ]; }; then
    echo "ERROR: Go $GO_VER is too old. Need Go 1.26+ (the SIMD experiment lives there)." >&2
    exit 1
fi

ARCH=$(go env GOARCH)
OS=$(go env GOOS)

# --- environment fingerprint ---

{
    echo "Date:       $(date -Iseconds 2>/dev/null || date)"
    echo "Hostname:   $(hostname)"
    echo "Kernel:     $(uname -srm)"
    echo "Go:         $(go version)"
    echo "GOOS/ARCH:  $OS/$ARCH"
    echo "Mode:       $mode (count=$COUNT, benchtime=$TIME)"
    echo ""
    echo "## CPU"
    if command -v lscpu >/dev/null 2>&1; then
        lscpu | grep -E '^(Architecture|Vendor ID|Model name|CPU\(s\)|Thread\(s\)|Core\(s\)|CPU MHz|CPU max MHz|L1d|L2|L3|Flags)' || true
    elif [ -r /proc/cpuinfo ]; then
        grep -m1 -E '^(model name|vendor_id|cpu MHz|cache size|flags)' /proc/cpuinfo || true
    else
        echo "(no lscpu and no /proc/cpuinfo — likely macOS or BSD)"
        if command -v sysctl >/dev/null 2>&1; then
            sysctl -n machdep.cpu.brand_string 2>/dev/null || true
            sysctl -n hw.ncpu 2>/dev/null || true
        fi
    fi
    echo ""
    echo "## CPU frequency scaling"
    if [ -d /sys/devices/system/cpu/cpu0/cpufreq ]; then
        gov=$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor 2>/dev/null || echo unknown)
        echo "Governor:   $gov"
        if [ "$gov" != "performance" ]; then
            echo ""
            echo "WARNING: governor is '$gov', not 'performance'."
            echo "Benchmark numbers will be noisier. To stabilize (temporary, root needed):"
            echo "  sudo cpupower frequency-set -g performance"
            echo "  # or:"
            echo "  for c in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor; do"
            echo "    echo performance | sudo tee \$c >/dev/null"
            echo "  done"
        fi
    else
        echo "(no /sys/devices/system/cpu/cpu0/cpufreq — not a Linux box, skipping)"
    fi
} | tee "$ENV_OUT"

# --- scalar run ---

echo ""
echo "==> running scalar bench (no GOEXPERIMENT)  [count=$COUNT, benchtime=$TIME]"
echo "    output: $SCALAR_OUT"
go test -bench=. -benchmem -count="$COUNT" -benchtime="$TIME" -run=^$ ./... \
    | tee "$SCALAR_OUT"

# --- SIMD run (linux/amd64 only) ---

if [ "$OS" = "linux" ] && [ "$ARCH" = "amd64" ]; then
    echo ""
    echo "==> running SIMD bench (GOEXPERIMENT=simd)  [count=$COUNT, benchtime=$TIME]"
    echo "    output: $SIMD_OUT"
    GOEXPERIMENT=simd go test -bench=. -benchmem -count="$COUNT" -benchtime="$TIME" -run=^$ ./... \
        | tee "$SIMD_OUT"

    # --- benchstat ---

    if ! command -v benchstat >/dev/null 2>&1; then
        echo ""
        echo "==> benchstat not in PATH; installing into \$(go env GOPATH)/bin"
        go install golang.org/x/perf/cmd/benchstat@latest
        export PATH="$(go env GOPATH)/bin:$PATH"
    fi

    echo ""
    echo "==> benchstat scalar vs simd"
    benchstat "$SCALAR_OUT" "$SIMD_OUT" | tee "$CMP_OUT"
else
    echo ""
    echo "==> skipping SIMD bench (need linux/amd64, got $OS/$ARCH)"
    : > "$SIMD_OUT"
    : > "$CMP_OUT"
fi

# --- summary markdown ---

{
    echo "# TurboSlice native bench — $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo ""
    echo "## Environment"
    echo '```'
    cat "$ENV_OUT"
    echo '```'
    if [ -s "$CMP_OUT" ]; then
        echo ""
        echo "## benchstat: scalar vs SIMD"
        echo '```'
        cat "$CMP_OUT"
        echo '```'
    fi
    echo ""
    echo "## Raw scalar"
    echo '```'
    cat "$SCALAR_OUT"
    echo '```'
    if [ -s "$SIMD_OUT" ]; then
        echo ""
        echo "## Raw SIMD"
        echo '```'
        cat "$SIMD_OUT"
        echo '```'
    fi
} > "$SUMMARY"

echo ""
echo "------------------------------------------------------------"
echo "DONE. Results in $OUT_DIR/"
echo ""
ls -lh "$OUT_DIR"
echo ""
echo "To send back, run:"
echo "  tar -czf bench-results.tgz $OUT_DIR/"
echo "...and email bench-results.tgz."
echo "(Or just send $SUMMARY — it has everything in one file.)"
