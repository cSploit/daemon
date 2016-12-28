FROM golang:alpine

MAINTAINER DeveloppSoft <developpsoft@gmail.com>

RUN apk --update add libpcap-dev git alpine-sdk
RUN rm -f /var/cache/apk/*

RUN go get -v github.com/cSploit/daemon

WORKDIR /go/src/github.com/cSploit/daemon

RUN go build -o ./cSploit

RUN mkdir /cSploit
RUN cp ./cSploit /cSploit/cSploit
RUN cp ./config.json /cSploit/config.json

WORKDIR /cSploit

# Add to path
RUN ln -s /cSploit/cSploit /bin/cSploit

# Build start.sh
RUN echo "export CSPLOIT_CONFIG='/cSploit/config.json'" > /start.sh
RUN echo "cSploit" >> /start.sh

RUN chmod +x /start.sh

CMD /start.sh

