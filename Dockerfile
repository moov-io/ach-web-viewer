# generated-from:4eb13d9dd6505954e3ad7dafa7ad1014b4f5710bce763c9b94bbb3e61ca04538 DO NOT REMOVE, DO UPDATE

FROM golang:1.20 as builder
WORKDIR /src
ARG VERSION

RUN apt-get update && apt-get install -y make gcc g++ ca-certificates

COPY . .

RUN VERSION=${VERSION} make build

FROM debian:stable-slim AS runtime
LABEL maintainer="Moov <oss@moov.io>"

WORKDIR /

RUN apt-get update && apt-get install -y ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/bin/ach-web-viewer /app/

ENV HTTP_PORT=8585
ENV HEALTH_PORT=9595

EXPOSE ${HTTP_PORT}/tcp
EXPOSE ${HEALTH_PORT}/tcp

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
	CMD curl -f http://localhost:${HEALTH_PORT}/live || exit 1

VOLUME [ "/data", "/configs" ]

ENTRYPOINT ["/app/ach-web-viewer"]
