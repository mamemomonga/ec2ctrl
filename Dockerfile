FROM golang:1.14.1-buster

RUN set -xe && \
	export DEBIAN_FRONTEND=noninteractive && \
	apt-get install -y --no-install-recommends \
		git-core make && \
	rm -rf /var/lib/apt/lists/*

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV GO111MODULE=on

ADD . /g/
WORKDIR /g
RUN set -xe && make multiarch
