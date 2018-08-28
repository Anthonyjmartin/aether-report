#!/bin/bash
APPVER="0.1.1"

# Build the binaries
docker build --rm -t aether-report-build -f Dockerfile .
docker run --rm -i -v $PWD:/go/src/gitlab.com/anthony.j.martin/aether-report aether-report-build sh -c "\
    echo GoLint; \
    golist=$(go list ./...);
    for i in ${golist}; do   golint $i; done; \
    echo; echo Go Vet; \
    go tool vet ./; \
    echo; echo Go Test and Coverage; \
    go test -coverprofile=coverage.out ./...; \
    go tool cover -html=coverage.out -o coverage.html"
#open ./coverage.html
