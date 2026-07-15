#!/bin/bash
set -e

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
BUILD_DIR="$ROOT_DIR/target"
VERSION="${1:-$(date +%Y%m%d)}"
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"
PKG_NAME="zonedesk-${VERSION}-${GOOS}-${GOARCH}"
STAGING="$(mktemp -d)/${PKG_NAME}"

echo "==> Building frontend..."
cd "$ROOT_DIR/frontend"
npm install
npm run build

echo "==> Building backend..."
cd "$ROOT_DIR/backend"
go build -o "$BUILD_DIR/zonedesk" ./cmd/server

echo "==> Copying frontend dist..."
rm -rf "$BUILD_DIR/dist"
mkdir -p "$BUILD_DIR/dist"
cp -r "$ROOT_DIR/frontend/dist"/* "$BUILD_DIR/dist/"

echo "==> Creating package structure..."
mkdir -p "$STAGING/app" "$STAGING/config" "$STAGING/data/cache" "$STAGING/logs" "$STAGING/web/dist"

echo "==> Copying binary..."
cp "$BUILD_DIR/zonedesk" "$STAGING/app/zonedesk"
chmod +x "$STAGING/app/zonedesk"

echo "==> Copying frontend dist..."
cp -r "$BUILD_DIR/dist"/* "$STAGING/web/dist/"

echo "==> Copying config..."
cp "$ROOT_DIR/scripts/deploy/config.yaml" "$STAGING/config/config.yaml"

echo "==> Copying start/stop scripts..."
cp "$ROOT_DIR/scripts/deploy/start.sh" "$STAGING/start.sh"
cp "$ROOT_DIR/scripts/deploy/stop.sh" "$STAGING/stop.sh"
chmod +x "$STAGING/start.sh" "$STAGING/stop.sh"

echo "==> Creating release archive..."
mkdir -p "$BUILD_DIR"
cd "$(dirname "$STAGING")"
tar czf "$BUILD_DIR/${PKG_NAME}.tar.gz" "$PKG_NAME"
rm -rf "$(dirname "$STAGING")"

echo ""
echo "==> Build completed."
echo "Binary:  $BUILD_DIR/zonedesk"
echo "Package: $BUILD_DIR/${PKG_NAME}.tar.gz"
echo "Size:    $(du -h "$BUILD_DIR/${PKG_NAME}.tar.gz" | cut -f1)"
echo ""
echo "Usage:"
echo "  tar xzf ${PKG_NAME}.tar.gz"
echo "  cd ${PKG_NAME}"
echo "  ./start.sh"
echo ""
echo "Visit: http://localhost:7070"
