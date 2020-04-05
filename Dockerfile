FROM golang:1.12 as builder

RUN mkdir -p $GOPATH/src/github.com/Visteras/vscale
COPY ./ $GOPATH/src/github.com/Visteras/vscale
WORKDIR $GOPATH/src/github.com/Visteras/vscale
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o .

FROM alpine
MAINTAINER Anatoliy Evladov <ae@visteras.ru>

RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /opt
COPY --from=builder /go/src/github.com/Visteras/vscale/vscale /opt/vscale
RUN chmod +x /opt/vscale

CMD ["/opt/vscale"]
