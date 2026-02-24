#!/bin/bash
# Build massdns from source (for WSL/Linux when the included binary fails due to glibc).
# Run from project root:  bash scripts/build-massdns-wsl.sh
# Requires: gcc (apt install build-essential if needed)

set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="/tmp/massdns-build-$$"
MASSDNS_REPO="https://github.com/blechschmidt/massdns.git"

echo "Cloning massdns..."
git clone --depth 1 "$MASSDNS_REPO" "$BUILD_DIR"
cd "$BUILD_DIR"

echo "Building massdns..."
mkdir -p bin
gcc -O3 -std=c11 -DHAVE_EPOLL -DHAVE_SYSINFO -Wall -fstack-protector-strong src/main.c -o bin/massdns

echo "Installing to project massdns/bin/massdns..."
cp bin/massdns "$PROJECT_ROOT/massdns/bin/massdns"
chmod +x "$PROJECT_ROOT/massdns/bin/massdns"

rm -rf "$BUILD_DIR"
echo "Done. Use: resodns resolve ... -b ./massdns/bin/massdns -r massdns/lists/resolvers.txt"
