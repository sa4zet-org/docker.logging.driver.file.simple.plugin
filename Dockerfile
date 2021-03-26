FROM golang:1.16.2-alpine3.12 as builder

WORKDIR github.com/sa4zet-org/docker.logging.driver.file.simple
COPY src/ .

RUN apk update \
  && apk add git \
  && go build \
    --ldflags "-extldflags -static" \
    -o /usr/bin/fileLoggingDriver

FROM alpine:3.12
COPY --from=builder /usr/bin/fileLoggingDriver /usr/bin/fileLoggingDriver
