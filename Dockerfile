FROM vibioh/scratch

ENTRYPOINT [ "/spurf" ]

ENV ZONEINFO /zoneinfo.zip
COPY zoneinfo.zip /zoneinfo.zip
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ARG VERSION
ENV VERSION=${VERSION}

ARG TARGETOS
ARG TARGETARCH

COPY release/spurf_${TARGETOS}_${TARGETARCH} /spurf
