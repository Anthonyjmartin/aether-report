#!/usr/bin/env bash
GOOS=linux go build -o build/linux/aether-report cmd/aether-report/main.go
rsync -av --progress build/* ~/Dropbox/gotesting/ --exclude .hold
