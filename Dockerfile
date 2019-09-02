FROM golang:1.12 as builder

WORKDIR /app
COPY . .

RUN make \
 && curl -q -sS -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && curl -q -sS -o /app/zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

ARG CODECOV_TOKEN
RUN curl -q -sS https://codecov.io/bash | bash

FROM scratch

ENV ZONEINFO zoneinfo.zip

HEALTHCHECK --retries=10 CMD [ "/spurf", "-c" ]
ENTRYPOINT [ "/spurf" ]

ARG APP_VERSION
ENV VERSION=${APP_VERSION}

COPY --from=builder /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/zoneinfo.zip /app/bin/spurf /
