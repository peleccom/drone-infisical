FROM alpine:3.17 as alpine
RUN apk add -U --no-cache ca-certificates

FROM golang:1.24-alpine as builder
RUN apk add --no-cache build-base

WORKDIR /src
COPY . .
RUN go build -o /bin/server .



FROM alpine:3.17
EXPOSE 3000

ENV GODEBUG netdns=go

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/server /bin/

ENTRYPOINT ["/bin/server"]
