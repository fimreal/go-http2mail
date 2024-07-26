FROM golang:latest AS builder
COPY . /app
RUN cd /app &&\
    make docker-build &&\
    ls -l bin 

# 下载证书
FROM alpine:latest AS ca
RUN apk --no-cache add ca-certificates

# 
FROM scratch
LABEL source.url="https://github.com/fimreal/go-http2mail"

COPY --from=ca /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/http2mail /http2mail

EXPOSE 8782
ENTRYPOINT [ "/http2mail" ]