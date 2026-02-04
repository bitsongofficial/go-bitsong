#!/usr/bin/env bash
set -euo pipefail

# This script creates release artifacts from reproducible builds.
# Run `make build-reproducible` first to generate the binaries.
#
# Usage: ./scripts/prep-release.sh [VERSION]
# If VERSION is not provided, it will be extracted from git tags.

VERSION="${1:-$(git describe --tags 2>/dev/null | sed 's/^v//' || echo "unknown")}"
BUILD_DIR="build"

echo "Preparing release artifacts for version: $VERSION"

# Check if binaries exist
if [[ ! -f "$BUILD_DIR/bitsongd-linux-amd64" ]] && [[ ! -f "$BUILD_DIR/bitsongd-linux-arm64" ]]; then
    echo "Error: No binaries found in $BUILD_DIR/"
    echo "Run 'make build-reproducible' first."
    exit 1
fi

# Create tarballs and checksums
for arch in amd64 arm64; do
    binary="$BUILD_DIR/bitsongd-linux-$arch"
    if [[ -f "$binary" ]]; then
        tarball="$BUILD_DIR/bitsongd-$VERSION-linux-$arch.tar.gz"
        echo "Creating $tarball..."
        tar -czvf "$tarball" -C "$BUILD_DIR" "bitsongd-linux-$arch"

        echo "Generating checksum..."
        sha256sum "$tarball" > "$tarball.sha256"
    fi
done

# Generate combined checksum file
echo "Generating combined checksum file..."
cat "$BUILD_DIR"/*.sha256 > "$BUILD_DIR/checksums.txt" 2>/dev/null || true

echo ""
echo "Release artifacts created in $BUILD_DIR/:"
ls -la "$BUILD_DIR"/*.tar.gz "$BUILD_DIR"/*.sha256 2>/dev/null || echo "No artifacts found"
