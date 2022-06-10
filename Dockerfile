FROM 1.17.11-alpine3.16 as builder

WORKDIR github.com/sa4zet-org/docker.logging.driver.file.simple
COPY src/ .

RUN apk update \
  && apk add git \
  && go build \
    --ldflags "-extldflags -static" \
    -o /usr/bin/fileLoggingDriver

FROM alpine:3.16
COPY --from=builder /usr/bin/fileLoggingDriver /usr/bin/fileLoggingDriver
