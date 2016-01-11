FROM alpine:latest
MAINTAINER Kai Hendry <hendry@iki.fi>

COPY . /go/src/github.com/kaihendry/lk
RUN apk add --no-cache go git make \
	&& cd /go/src/github.com/kaihendry/lk \
	&& export GOPATH=/go \
	&& go get \
	&& make \
	&& mv lk /bin/lk \
	&& rm -rf /go

COPY i /srv

EXPOSE 3000
ENTRYPOINT ["/bin/lk", "--port", "3000"]
CMD ["/srv"]
