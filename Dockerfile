FROM golang:1.14 as builder
ADD . /src/go-http2mail
RUN cd /src/go-http2mail &&\
    go mod download &&\
    make docker

FROM scratch
# ADD go-http2mail /
COPY --from=builder /src/go-http2mail/bin/go-http2mail-docker /go-http2mail
EXPOSE 5000
ENTRYPOINT ["/go-http2mail"]
