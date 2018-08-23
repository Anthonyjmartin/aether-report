FROM golang:1.10-alpine
WORKDIR /go/src/gitlab.com/anthony.j.martin/aether-report
ADD ./ /go/src/gitlab.com/anthony.j.martin/aether-report
ENV GOOS=linux
