FROM golang:1.20-alpine3.17 AS builder

ENV GOBIN=/build
ENV CGO_ENABLED=0
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY ./ ./
RUN go build ./...

FROM alpine:3.17
COPY --from=builder /build/sonapi-bot /usr/local/bin/sonapi-bot

RUN apk add --no-cache ca-certificates && update-ca-certificates && apk del busybox

ENTRYPOINT ["/usr/local/bin/sonapi-bot"]
