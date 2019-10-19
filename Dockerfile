FROM golang:1.13 as builder

WORKDIR /app
COPY . .

RUN make

ARG CODECOV_TOKEN
RUN curl -q -sS https://codecov.io/bash | bash

FROM alpine as fetcher

WORKDIR /app

RUN apk --update add curl \
 && curl -q -sS -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && curl -q -sS -o /app/zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

FROM scratch

ENV ZONEINFO zoneinfo.zip

HEALTHCHECK --retries=10 CMD [ "/spurf", "-c" ]
ENTRYPOINT [ "/spurf" ]

ARG APP_VERSION
ENV VERSION=${APP_VERSION}

COPY --from=fetcher /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=fetcher /app/zoneinfo.zip /
COPY --from=builder /app/bin/spurf /
