#!/bin/bash
ZONE_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_FILE="$ZONE_DIR/zonedesk.pid"
LOG_FILE="$ZONE_DIR/logs/zonedesk.log"

if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
  echo "ZoneDesk is already running (PID: $(cat "$PID_FILE"))"
  exit 1
fi

mkdir -p "$ZONE_DIR/logs" "$ZONE_DIR/data/cache"

cd "$ZONE_DIR"
nohup ./app/zonedesk --config ./config/config.yaml >> "$LOG_FILE" 2>&1 &
echo $! > "$PID_FILE"

sleep 1
if kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
  echo "ZoneDesk started (PID: $(cat "$PID_FILE"))"
  echo "Log: $LOG_FILE"
else
  echo "Failed to start ZoneDesk. Check log: $LOG_FILE"
  rm -f "$PID_FILE"
  exit 1
fi
