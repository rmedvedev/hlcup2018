FROM golang:latest AS builder

WORKDIR $GOPATH/src/app/

RUN go get -u github.com/valyala/fasthttp
RUN go get -u github.com/google/btree

COPY src $GOPATH/src/app/

RUN go build -o /go/bin/highloadcup-service

FROM centos:7

COPY --from=builder /go/bin/highloadcup-service /go/bin/highloadcup-service

EXPOSE 80

CMD /go/bin/highloadcup-service
