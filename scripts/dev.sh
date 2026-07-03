#!/bin/bash
set -e

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)

echo "==> Starting backend..."
cd "$ROOT_DIR/backend"
mkdir -p data/cache logs
go run ./cmd/server --config "$ROOT_DIR/deploy/config.yaml" &
BACKEND_PID=$!

echo "==> Starting frontend dev server..."
cd "$ROOT_DIR/frontend"
npm run dev &
FRONTEND_PID=$!

trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null" EXIT

echo "Backend PID: $BACKEND_PID"
echo "Frontend PID: $FRONTEND_PID"
echo "Frontend: http://localhost:5173"
echo "Backend:  http://localhost:7070"

wait
