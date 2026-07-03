#!/bin/bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="${1:-$(date +%Y%m%d)}"
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"
OUTPUT_DIR="$ROOT_DIR/build"
PKG_NAME="zonedesk-${VERSION}-${GOOS}-${GOARCH}"
STAGING="$(mktemp -d)/${PKG_NAME}"

if [ ! -f "$ROOT_DIR/build/zonedesk" ] || [ ! -d "$ROOT_DIR/build/dist" ]; then
  echo "==> Build artifacts not found. Running build.sh..."
  "$ROOT_DIR/scripts/build.sh"
fi

echo "==> Creating package structure..."
mkdir -p "$STAGING/app" "$STAGING/config" "$STAGING/data/cache" "$STAGING/logs" "$STAGING/web/dist"

echo "==> Copying binary..."
cp "$ROOT_DIR/build/zonedesk" "$STAGING/app/zonedesk"
chmod +x "$STAGING/app/zonedesk"

echo "==> Copying frontend dist..."
cp -r "$ROOT_DIR/build/dist"/* "$STAGING/web/dist/"

echo "==> Copying config..."
cp "$ROOT_DIR/deploy/config.yaml" "$STAGING/config/config.yaml"

echo "==> Copying start/stop scripts..."
cp "$ROOT_DIR/deploy/start.sh" "$STAGING/start.sh"
cp "$ROOT_DIR/deploy/stop.sh" "$STAGING/stop.sh"
chmod +x "$STAGING/start.sh" "$STAGING/stop.sh"

echo "==> Creating release archive..."
mkdir -p "$OUTPUT_DIR"
cd "$(dirname "$STAGING")"
tar czf "$OUTPUT_DIR/${PKG_NAME}.tar.gz" "$PKG_NAME"
rm -rf "$(dirname "$STAGING")"

echo ""
echo "Package created: $OUTPUT_DIR/${PKG_NAME}.tar.gz"
echo "Size: $(du -h "$OUTPUT_DIR/${PKG_NAME}.tar.gz" | cut -f1)"
echo ""
echo "Usage:"
echo "  tar xzf ${PKG_NAME}.tar.gz"
echo "  cd ${PKG_NAME}"
echo "  ./start.sh"
echo ""
echo "Visit: http://localhost:7070"
