#!/bin/bash
exec > /tmp/post-create.log 2>&1
set -x

echo "=== WA Farm Setup ==="
date

cd /workspaces/wa-farm

go version

# Build wabot (WhatsApp bridge)
echo "--- Building wabot ---"
cd wabot
go mod tidy 2>&1
go build -o /workspaces/wa-farm/wabot . 2>&1
WABOT_OK=$?
cd /workspaces/wa-farm

# Build main server
echo "--- Building server ---"
go mod tidy 2>&1
go build -o server server.go 2>&1

# Start
echo "--- Starting server ---"
nohup /workspaces/wa-farm/server > /tmp/server.log 2>&1 &
sleep 2
curl -s http://localhost:8080/status || echo "Server check failed"
echo "=== Done ==="
