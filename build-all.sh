#!/bin/bash
echo Building linux x64
env GOOS=linux GOARCH=amd64 go build -o ./build/webhook_linux_x64

echo Building linux arm64
env GOOS=linux GOARCH=arm64 go build -o ./build/webhook_linux_arm64

echo Building darwin amd64
env GOOS=darwin GOARCH=amd64 go build -o ./build/webhook_darwin_amd64

echo Building darwin arm64
env GOOS=darwin GOARCH=arm64 go build -o ./build/webhook_darwin_arm64

echo Building windows amd64
env GOOS=windows GOARCH=amd64 go build -o ./build/webhook_windows_amd64

echo Building windows arm64
env GOOS=windows GOARCH=arm64 go build -o ./build/webhook_windows_arm64

echo Done ✅