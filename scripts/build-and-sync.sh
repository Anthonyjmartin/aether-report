#!/usr/bin/env bash

# Build the binaries
docker build --rm -t aether-report-build -f Dockerfile .
docker run --rm -i -v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report aether-report-build go build -a -v -o build/linux/amd64/aether-report -ldflags="-X main.version=0.1.0"  cmd/aether-report/main.go #./build.sh
docker rmi aether-report-build

# Build the RPM
docker build --rm -t aether-report-rpm -f Dockerfile.RPM .
docker run --rm -i -v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report aether-report-rpm go-bin-rpm generate -a amd64 -o build/linux/amd64/aether-report.rpm --version 0.1.0
docker rmi aether-report-rpm

#GOOS=linux go build -a -o build/linux/amd64/aether-report -ldflags="-X main.version=0.1.0"  cmd/aether-report/main.go
rsync -av --progress build/* ~/Dropbox/gotesting/ --exclude .hold
