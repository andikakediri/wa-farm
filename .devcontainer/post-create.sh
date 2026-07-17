#!/bin/bash
exec > /tmp/post-create.log 2>&1
set -x

echo "================================================"
echo "  WA Farm - Full Setup"
echo "================================================"
date

cd /workspaces/wa-farm

# Step 1: Initialize Go modules
echo "[1/5] Initializing Go modules..."
go mod tidy 2>&1
echo "go mod tidy done"

# Step 2: Build the wabot (WhatsApp protocol bridge)
echo "[2/5] Building wabot (whatsmeow bridge)..."
cd cmd/wabot
go build -o /workspaces/wa-farm/wabot . 2>&1
WABOT_EXIT=$?
echo "wabot build exit: $WABOT_EXIT"
cd /workspaces/wa-farm

# Step 3: Build main server
echo "[3/5] Building server..."
go build -o server server.go 2>&1
echo "server build exit: $?"

# Step 4: Check if wabot built successfully
echo "[4/5] Checking binaries..."
if [ -f wabot ]; then
    echo "wabot binary exists ($(ls -la wabot | awk '{print $5}') bytes)"
else
    echo "wabot binary NOT built - server will use simulation mode"
fi
ls -la server

# Step 5: Start server
echo "[5/5] Starting server..."
pkill -f "^./server" 2>/dev/null || true
nohup ./server > /tmp/server.log 2>&1 &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"
sleep 2

# Check if running
curl -s http://localhost:8080/status || echo "Server not responding yet"

echo ""
echo "================================================"
echo "  Server: PID $SERVER_PID"
echo "  Mode: $(test -f wabot && echo 'REAL (whatsmeow)' || echo 'SIMULATION')"
echo "  Logs: /tmp/server.log"
echo "================================================"
