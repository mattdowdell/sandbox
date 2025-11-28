# https://docs.docker.com/reference/dockerfile/
#
# mirror.gcr.io caches of popular docker hub images, but does not add rate limiting.
# See https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images.

# ------------
# Build target
# ------------

FROM --platform=$BUILDPLATFORM mirror.gcr.io/golang:1.25-bookworm@sha256:e17419604b6d1f9bc245694425f0ec9b1b53685c80850900a376fb10cb0f70cb AS build

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

FROM gcr.io/distroless/static-debian12:nonroot@sha256:e8a4044e0b4ae4257efa45fc026c0bc30ad320d43bd4c1a7d5271bd241e386d0 AS runtime

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
COPY --from=build /go/bin/example-health /example-health
