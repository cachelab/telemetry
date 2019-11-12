FROM golang:1.13.4-alpine

MAINTAINER Cache Lab <hello@cachelab.co>

COPY telemetry /bin/telemetry

USER nobody

ENTRYPOINT ["/bin/telemetry"]
