# Base image: https://hub.docker.com/_/golang/
FROM golang:1.10
MAINTAINER Anthony Martin <anthony.j.martin142@gmail.com>

ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH
RUN go get -u github.com/golang/lint/golint

RUN wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add -; \
    echo "deb http://apt.llvm.org/stretch/ llvm-toolchain-stretch-5.0 main" | tee -a /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
    clang-5.0 \
    rpm \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*; \
    curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh | GH=mh-cbon/go-bin-rpm sh -xe; \
    curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh | GH=mh-cbon/changelog sh -xe

# Set Clang as default CC
ENV set_clang /etc/profile.d/set-clang-cc.sh
RUN echo "export CC=clang-5.0" | tee -a ${set_clang} && chmod a+x ${set_clang}
