#!/bin/bash
# Run TurboSlice benchmarks on AMD64 via Docker (with QEMU emulation on ARM64 hosts).
# Usage: ./scripts/bench-amd64.sh
#
# Note: On ARM64 hosts this uses QEMU emulation so numbers will be ~5-10x slower
# than native AMD64. The RELATIVE comparison (SIMD vs scalar) is still valid.
# For accurate absolute numbers, run on a native AMD64 machine.

set -e

cd "$(dirname "$0")/.."

echo "Building AMD64 benchmark container..."
docker build --platform linux/amd64 -f Dockerfile.bench -t turboslice-bench . 2>&1 | tail -3

echo ""
echo "Running benchmarks (AMD64 via QEMU)..."
echo "This may take a few minutes."
echo ""

docker run --rm --platform linux/amd64 turboslice-bench
