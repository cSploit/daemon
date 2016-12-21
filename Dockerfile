FROM golang:alpine

MAINTAINER DeveloppSoft <developpsoft@gmail.com>

RUN apk --update add libpcap-dev git alpine-sdk
RUN rm -f /var/cache/apk/*

RUN go get github.com/cSploit/daemon
RUN go install github.com/cSploit/daemon

ENTRYPOINT /go/bin/daemon

