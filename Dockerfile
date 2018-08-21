FROM golang:1.10
WORKDIR /go/src/gitlab.com/anthony.j.martin/aether-report
ADD ./ /go/src/gitlab.com/anthony.j.martin/aether-report
ENV GOOS=linux
