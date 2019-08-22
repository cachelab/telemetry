FROM golang:1.12.9-alpine

MAINTAINER Cache Lab <hello@cachelab.co>

COPY telemetry /bin/telemetry

USER nobody

ENTRYPOINT ["/bin/telemetry"]
