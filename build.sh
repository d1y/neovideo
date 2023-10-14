#!/bin/bash

pushd frontend
cd admin
bun i
bun run build
cd ../appvideo
bun i
bun run build
popd

go build -o appvideo.exe .