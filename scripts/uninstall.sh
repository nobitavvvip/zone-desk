#!/bin/bash
set -e

ZONE_ROOT="${1:-/srv/zonedesk}"

echo "WARNING: This will remove ZoneDesk from: $ZONE_ROOT"
read -p "Are you sure? (yes/NO): " confirm
if [ "$confirm" != "yes" ]; then
  echo "Aborted."
  exit 1
fi

echo "==> Stopping ZoneDesk..."
if [ -f "$ZONE_ROOT/stop.sh" ]; then
  "$ZONE_ROOT/stop.sh" 2>/dev/null || true
fi

echo "==> Removing ZoneDesk files..."
rm -rf "$ZONE_ROOT"

echo "ZoneDesk has been uninstalled."
echo "Tip: pass a custom path to uninstall from elsewhere:"
echo "  $0 /opt/my-zonedesk"
