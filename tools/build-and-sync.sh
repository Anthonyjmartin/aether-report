#!/usr/bin/env bash

APPVER="0.1.1"

# Build the binaries
docker build --rm -t aether-report-build -f Dockerfile .
docker run --rm -i -v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report aether-report-build go build -a -v -o build/linux/amd64/aether-report -ldflags="-X main.version=${APPVER}"  cmd/aether-report/main.go
#docker rmi aether-report-build

# Build the RPM
docker build --rm -t aether-report-rpm -f Dockerfile.RPM .
docker run --rm -i -v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report aether-report-rpm go-bin-rpm generate -a amd64 -o build/linux/amd64/aether-report-${APPVER}.rpm --version ${APPVER}
#docker rmi aether-report-rpm

#GOOS=linux go build -a -o build/linux/amd64/aether-report -ldflags="-X main.version=${APPVER}"  cmd/aether-report/main.go
rsync -av --progress build/* ~/Dropbox/gotesting/ --exclude .hold
a
