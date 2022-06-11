FROM golang:1.17.11-alpine3.16 as builder
RUN mkdir /tmp/build
WORKDIR /tmp/build
COPY src/ .
RUN apk update \
  && apk add git \
  && CGO_ENABLED=0 GOOS=linux go build --ldflags '-extldflags "-static -fno-PIC -O2"' -a -v -o /usr/bin/fileLoggingDriver
FROM alpine:3.16
COPY --from=builder /usr/bin/fileLoggingDriver /usr/bin/fileLoggingDriver
