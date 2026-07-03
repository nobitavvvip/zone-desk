#!/bin/bash
set -e

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
BUILD_DIR="$ROOT_DIR/build"

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

echo "==> Build completed. Output in $BUILD_DIR/"
