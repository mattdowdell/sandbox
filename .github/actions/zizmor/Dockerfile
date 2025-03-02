# https://docs.docker.com/reference/dockerfile/

FROM ghcr.io/woodruffw/zizmor:1.4.1@sha256:dbd8fffd06c51b0ea80dff6115547a0b4fe97fc67785d82f198ecf36d8a32ebe

COPY --chmod=555 entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
