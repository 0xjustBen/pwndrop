#!/bin/bash
echo Building...
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/pwndrop main.go
echo Done.