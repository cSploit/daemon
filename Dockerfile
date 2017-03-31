FROM golang:alpine

MAINTAINER DeveloppSoft <developpsoft@gmail.com>

RUN apk --update add libpcap-dev git alpine-sdk
RUN rm -f /var/cache/apk/*

RUN mkdir -p /go/src/github.com/cSploit/daemon
ADD . /go/src/github.com/cSploit/daemon

WORKDIR /go/src/github.com/cSploit/daemon

RUN go get -v ./...
RUN go build

CMD ./daemon

