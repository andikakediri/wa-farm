#!/bin/bash
exec > /tmp/post-create.log 2>&1

echo "=== WA Farm Setup ==="
date
go version
cd /workspaces/wa-farm

# 1. Main server (should always work)
echo "--- Building server ---"
go mod tidy 2>&1
go build -o server server.go 2>&1
echo "server build: $?"

# 2. Try to build wabot (optional, might fail without new Go)
echo "--- Trying wabot ---"
cd wabot
go mod tidy 2>&1
go build -o /workspaces/wa-farm/wabot . 2>&1
WABOT=$?
echo "wabot build: $WABOT"
cd /workspaces/wa-farm

# 3. Start
echo "--- Starting ---"
nohup ./server > /tmp/server.log 2>&1 &
sleep 2
curl -s http://localhost:8080/status
echo ""
echo "=== Done (wabot: $WABOT) ==="
