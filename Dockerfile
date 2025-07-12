# https://docs.docker.com/reference/dockerfile/
#
# mirror.gcr.io caches of popular docker hub images, but does not add rate limiting.
# See https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images.

# ------------
# Build target
# ------------

FROM --platform=$BUILDPLATFORM mirror.gcr.io/golang:1.24-bookworm@sha256:69adc37c19ac6ef724b561b0dc675b27d8c719dfe848db7dd1092a7c9ac24bc6 AS build

WORKDIR /go/src

ARG TARGETOS TARGETARCH
ARG SOURCE_DATE_EPOCH=0
ARG GO_BUILD_ARGS=
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=. \
    set -eux; \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOTOOLCHAIN=local \
    go build ${GO_BUILD_ARGS} -trimpath -ldflags="-buildid= -s -w" -o /go/bin/ ./cmd/...; \
    touch --date=@${SOURCE_DATE_EPOCH} /go/bin/*;

# --------------
# Runtime target
# --------------

FROM gcr.io/distroless/static-debian12:nonroot@sha256:627d6c5a23ad24e6bdff827f16c7b60e0289029b0c79e9f7ccd54ae3279fb45f AS runtime

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
COPY --from=build /go/bin/example-health /example-health
