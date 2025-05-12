# https://docs.docker.com/reference/dockerfile/

FROM ghcr.io/woodruffw/zizmor:1.7.0@sha256:3d53be4e0f48952112aa4ca00f45e724b70598082e7202c5c412a6779ca38134

COPY --chmod=555 entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
