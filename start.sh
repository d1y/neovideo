#!/bin/bash
chmod +x /app/av.exe
nohup /app/av.exe --conf=/app/config.yaml &
nginx -g "daemon off;"