#!/bin/bash
cd /workspaces/wa-farm
go mod tidy 2>&1
go build -o server server.go 2>&1
nohup ./server > /tmp/server.log 2>&1 &
sleep 3
curl -s http://localhost:8080/status
