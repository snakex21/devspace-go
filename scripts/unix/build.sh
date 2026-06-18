#!/bin/bash
# Build DevSpace for all platforms
# Works on: Linux, macOS, WSL, Git Bash (Windows)
# Run from project root: ./scripts/unix/build.sh

set -e
cd "$(dirname "$0")/../.."

PLATFORMS=(
    "windows:amd64:.exe"
    "linux:amd64:"
    "darwin:amd64:"
    "darwin:arm64:"
)

DIRS=("windows" "linux" "macos-intel" "macos-mchip")

echo "=== DevSpace Build All ==="
echo ""

rm -rf build
mkdir -p build

for i in "${!PLATFORMS[@]}"; do
    IFS=':' read -r os arch ext <<< "${PLATFORMS[$i]}"
    dir="${DIRS[$i]}"
    mkdir -p "build/$dir"

    echo "  [$dir] devspace (serwer)..."
    GOOS=$os GOARCH=$arch go build -o "build/$dir/devspace$ext" ./cmd/devspace/
    echo "        -> build/$dir/devspace$ext"

    echo "  [$dir] devspace-gui (konfigurator)..."
    GOOS=$os GOARCH=$arch go build -o "build/$dir/devspace-gui$ext" ./cmd/devspace-gui/ 2>/dev/null || \
        echo "        (GUI niedostępne przy cross-kompilacji — skompiluj natywnie)"

    echo ""
done

echo "=== Done! ==="
echo ""
for dir in "${DIRS[@]}"; do
    echo "--- $dir ---"
    ls -lh "build/$dir/" 2>/dev/null | tail -n +2 | awk '{print "  " $NF " (" $5 ")"}'
    echo ""
done
