#!/bin/bash
exec > /tmp/post-create.log 2>&1

echo "================================================"
echo "  WA Farm - Setup (Simulation Mode)"
echo "================================================"
date

cd /workspaces/wa-farm

echo "[1/3] go mod tidy..."
go mod tidy 2>&1

echo "[2/3] Building server..."
go build -o server server.go 2>&1

echo "[3/3] Starting server..."
nohup ./server > /tmp/server.log 2>&1 &
sleep 2

echo "Status: $(curl -s http://localhost:8080/status || echo 'FAIL')"
echo "Ready!"
