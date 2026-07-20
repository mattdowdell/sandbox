# https://docs.docker.com/reference/dockerfile/
#
# mirror.gcr.io caches of popular docker hub images, but does not add rate limiting.
# See https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images.

# ------------
# Build target
# ------------

FROM --platform=$BUILDPLATFORM mirror.gcr.io/golang:1.26-trixie@sha256:4ee9ffa999b4583ce281939cdff828763083610292f252279a0cee77473bd9a7 AS build

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

# ------------
# Debug target
# ------------

FROM gcr.io/distroless/static-debian13:debug-nonroot@sha256:e6e0c6315dec9279f0528dbf9dd0ca6ff5a2ecb8edf34347966828168ddb9545 AS debug

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
COPY --from=build /go/bin/example-health /example-health

# --------------
# Runtime target
# --------------

FROM gcr.io/distroless/static-debian13:nonroot@sha256:d29e660cc75a5b6b1334e03c5c81ccf9bc0884a002c6000dbf0fb96034814478 AS runtime

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
COPY --from=build /go/bin/example-health /example-health
