#!/bin/bash
set -e

ROOT_DIR=$(cd "$(dirname "$0")/" && pwd)
ZONE_ROOT="${1:-/srv/zonedesk}"

echo "当前执行目录: ${ROOT_DIR}"
echo "ZoneDesk安装目录: ${ZONE_ROOT}"

if [ ! -f "$ROOT_DIR/build/zonedesk" ]; then
  echo "Error: build/zonedesk not found. Run build.sh first."
  exit 1
fi

echo "==> Creating directories..."
mkdir -p "$ZONE_ROOT/app" "$ZONE_ROOT/config" "$ZONE_ROOT/data/cache" "$ZONE_ROOT/logs" "$ZONE_ROOT/web/dist"

echo "==> Installing binary..."
cp "$ROOT_DIR/build/zonedesk" "$ZONE_ROOT/app/zonedesk"
chmod +x "$ZONE_ROOT/app/zonedesk"

echo "==> Installing frontend dist..."
rm -rf "$ZONE_ROOT/web/dist"/*
cp -r "$ROOT_DIR/build/dist"/* "$ZONE_ROOT/web/dist/"

echo "==> Installing config..."
if [ ! -f "$ZONE_ROOT/config/config.yaml" ]; then
  cp "$ROOT_DIR/deploy/config.yaml" "$ZONE_ROOT/config/config.yaml"
  echo "  -> Created default config at $ZONE_ROOT/config/config.yaml"
fi

echo "==> Installing start/stop scripts..."
cp "$ROOT_DIR/deploy/start.sh" "$ZONE_ROOT/start.sh"
cp "$ROOT_DIR/deploy/stop.sh" "$ZONE_ROOT/stop.sh"
chmod +x "$ZONE_ROOT/start.sh" "$ZONE_ROOT/stop.sh"

echo ""
echo "ZoneDesk installed to: $ZONE_ROOT"
echo ""
echo "Start: $ZONE_ROOT/start.sh"
echo "Stop:  $ZONE_ROOT/stop.sh"
echo ""
echo "Access at: http://$(hostname -I | awk '{print $1}'):7070"
echo ""
echo "Tip: pass a custom path to install elsewhere:"
echo "  $0 /opt/my-zonedesk"
