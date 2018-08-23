#!/bin/bash
APPVER="0.1.1"

# Build the binaries
docker build --rm -t aether-report-build -f Dockerfile .
docker run --rm -i -v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report aether-report-build go test -coverprofile=coverage.out ./... # internal/pkg/hardware_check
docker run --rm -i -v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report aether-report-build go tool cover -html=coverage.out -o coverage.html
open ./coverage.html
