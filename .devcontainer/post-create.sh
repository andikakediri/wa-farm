#!/bin/bash

echo "================================================"
echo "  WA Farm - Auto Setup"
echo "================================================"
echo ""

cd /workspaces/wa-farm

# Step 1: Initialize Go modules
echo "[1/4] Initializing Go modules..."
go mod tidy 2>&1

echo ""
echo "[2/4] Building whatsmeow (WhatsApp protocol library)..."
cd whatsmeow
go build -o /workspaces/wa-farm/wabot ./example/ 2>&1
cd /workspaces/wa-farm

echo ""
echo "[3/4] Building server..."
go build -o server server.go 2>&1

echo ""
echo "[4/4] Starting server..."
nohup ./server > server.log 2>&1 &
SERVER_PID=$!
echo "Server running on PID: $SERVER_PID"

echo ""
echo "================================================"
echo "  ✅ SERVER RUNNING!"
echo "  Port: 8080"
echo "  PID: $SERVER_PID"
echo ""
echo "  Akses landing page:"
echo "  https://[CODESPACE_NAME]-8080.preview.app.github.dev"
echo ""
echo "  Check status:"
echo "  curl http://localhost:8080/status"
echo "================================================"
