#!/bin/bash
ZONE_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_FILE="$ZONE_DIR/zonedesk.pid"

if [ ! -f "$PID_FILE" ]; then
  echo "ZoneDesk is not running (no PID file)"
  exit 0
fi

PID=$(cat "$PID_FILE")
if kill -0 "$PID" 2>/dev/null; then
  echo "Stopping ZoneDesk (PID: $PID)..."
  kill "$PID"
  sleep 1
  if kill -0 "$PID" 2>/dev/null; then
    echo "Waiting for graceful shutdown..."
    sleep 2
    kill -9 "$PID" 2>/dev/null || true
  fi
  echo "ZoneDesk stopped."
else
  echo "ZoneDesk is not running (stale PID)"
fi

rm -f "$PID_FILE"
