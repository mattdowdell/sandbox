# https://docs.docker.com/reference/dockerfile/

FROM ghcr.io/woodruffw/zizmor:1.5.1@sha256:59c0a546bcfe4d6c040fe1532ddc2f6ccc63e59ebd6bc4f8f4ec741d1b77899d

COPY --chmod=555 entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
