#!/bin/bash
set -e

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
ZONE_ROOT="${1:-/srv/zonedesk}"

if [ ! -f "$ROOT_DIR/build/zonedesk" ]; then
  echo "Error: build/zonedesk not found. Run build.sh first."
  exit 1
fi

echo "==> Stopping ZoneDesk..."
"$ZONE_ROOT/stop.sh" 2>/dev/null || true

echo "==> Backing up current version..."
BACKUP_DIR="/tmp/zonedesk-backup-$(date +%Y%m%d%H%M%S)"
mkdir -p "$BACKUP_DIR"
cp "$ZONE_ROOT/app/zonedesk" "$BACKUP_DIR/zonedesk" 2>/dev/null || true

echo "==> Updating binary..."
cp "$ROOT_DIR/build/zonedesk" "$ZONE_ROOT/app/zonedesk"
chmod +x "$ZONE_ROOT/app/zonedesk"

echo "==> Updating frontend..."
rm -rf "$ZONE_ROOT/web/dist"/*
cp -r "$ROOT_DIR/build/dist"/* "$ZONE_ROOT/web/dist/"

echo "==> Starting ZoneDesk..."
"$ZONE_ROOT/start.sh"

echo "==> Upgrade completed. Backup at $BACKUP_DIR"
