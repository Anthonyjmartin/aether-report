FROM golang:1.10-alpine
WORKDIR /go/src/gitlab.com/anthony.j.martin/aether-report
ADD ./ /go/src/gitlab.com/anthony.j.martin/aether-report
ENV GOOS=linux
RUN apk add git; \
    go get github.com/golang/dep/cmd/dep; \
    go get github.com/golang/lint/golint; \
    go get github.com/tebeka/go2xunit
