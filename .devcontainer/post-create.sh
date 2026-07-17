#!/bin/bash
exec > /tmp/post-create.log 2>&1
set -x

echo "================================================"
echo "  WA Farm - Auto Setup"
echo "================================================"
date

cd /workspaces/wa-farm

# Step 1: Initialize Go modules
echo "[1/4] Initializing Go modules..."
go mod tidy 2>&1
echo "go mod tidy done"

echo ""
echo "[2/4] Building server (simulation mode - no wabot)..."
go build -o server server.go 2>&1
echo "Build exit code: $?"

echo ""
echo "[3/4] Starting server..."
pkill -f "server$" 2>/dev/null || true
nohup ./server > /tmp/server.log 2>&1 &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"
sleep 2

echo ""
echo "[4/4] Checking server..."
curl -s http://localhost:8080/status || echo "Server not responding yet"

echo ""
echo "================================================"
echo "  Server PID: $SERVER_PID"
echo "  Log: /tmp/server.log"
echo "  Post-create log: /tmp/post-create.log"
echo "================================================"
