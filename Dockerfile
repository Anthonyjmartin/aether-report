# Base image: https://hub.docker.com/_/golang/
FROM golang:1.10
MAINTAINER Anthony Martin <anthony.j.martin142@gmail.com>

# Install golint
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH
RUN go get -u github.com/golang/lint/golint

# Add apt key for LLVM repository
RUN wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add -

# Add LLVM apt repository
RUN echo "deb http://apt.llvm.org/stretch/ llvm-toolchain-stretch-5.0 main" | tee -a /etc/apt/sources.list

# Install clang from LLVM repository
RUN apt-get update && apt-get install -y --no-install-recommends \
    clang-5.0 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Set Clang as default CC
ENV set_clang /etc/profile.d/set-clang-cc.sh
RUN echo "export CC=clang-5.0" | tee -a ${set_clang} && chmod a+x ${set_clang}

#FROM golang:1.10
#WORKDIR /go/src/gitlab.com/anthony.j.martin/aether-report
#ADD ./ /go/src/gitlab.com/anthony.j.martin/aether-report
#ENV GOOS=linux version=""
#RUN go get github.com/golang/dep/cmd/dep; \
#    go get github.com/golang/lint/golint; \
#    go get github.com/tebeka/go2xunit
