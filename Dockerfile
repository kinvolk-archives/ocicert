FROM golang:alpine

RUN apk update && apk add --no-cache go git make

RUN export GOPATH="$(mktemp -d)" \
	&& export CGO_ENABLED=0 \
	&& export OCICERT_REGISTRY=127.0.0.1:5000/busybox:latest \
	&& export OCICERT_LOCALREG=1 \
	&& git clone https://github.com/docker/distribution "$GOPATH/src/github.com/docker/distribution" \
	&& cd "$GOPATH/src/github.com/docker/distribution" \
	&& go build -o /usr/local/bin/registry github.com/docker/distribution/cmd/registry \
	&& git clone https://github.com/kinvolk/ocicert "$GOPATH/src/github.com/kinvolk/ocicert" \
	&& cd "$GOPATH/src/github.com/kinvolk/ocicert" \
	&& git checkout dongsu/poc-integration-test \
	&& make test \
	&& rm -rf "$GOPATH"

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /go/src/github.com/kinvolk/ocicert
COPY . /go/src/github.com/kinvolk/ocicert
