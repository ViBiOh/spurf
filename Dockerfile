FROM vibioh/scratch

ENV ZONEINFO zoneinfo.zip

HEALTHCHECK --retries=10 CMD [ "/spurf", "-c" ]
ENTRYPOINT [ "/spurf" ]

ARG VERSION
ENV VERSION=${VERSION}

ARG TARGETOS
ARG TARGETARCH

COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY zoneinfo.zip /
COPY release/spurf_${TARGETOS}_${TARGETARCH} /spurf
