FROM scratch

ARG TARGETPLATFORM

LABEL org.opencontainers.image.source="https://github.com/e2engine/cli"
LABEL org.opencontainers.image.description="Declarative end-to-end testing and service mocking for distributed systems."
LABEL org.opencontainers.image.licenses="Apache-2.0"

COPY ${TARGETPLATFORM}/e2engine /usr/bin/e2engine

ENTRYPOINT ["/usr/bin/e2engine"]