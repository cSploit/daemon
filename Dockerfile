FROM golang:alpine

MAINTAINER DeveloppSoft <developpsoft@gmail.com>

RUN apk --update add libpcap-dev git
RUN rm -f /var/cache/apk/*

WORKDIR /opt
ADD ./* cSploit/
WORKDIR /opt/cSploit

RUN go get
RUN go build -o ./cSploit
RUN ln -s cSploit /usr/local/bin

CMD ./cSploit

