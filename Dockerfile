# https://docs.docker.com/reference/dockerfile/
#
# mirror.gcr.io caches of popular docker hub images, but does not add rate limiting.
# See https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images.

# ------------
# Build target
# ------------

FROM --platform=$BUILDPLATFORM mirror.gcr.io/golang:1.26-trixie@sha256:d08bf3ed2bd263088ca8e23fefaf10f1b71769f6932f0a4017ba28d2a5baf001 AS build

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

FROM gcr.io/distroless/static-debian13:debug-nonroot@sha256:cca906a21ecac7923195567342dc58bbfeb314912bf14f9280b4bb418fd4be6d AS debug

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
COPY --from=build /go/bin/example-health /example-health

# --------------
# Runtime target
# --------------

FROM gcr.io/distroless/static-debian13:nonroot@sha256:e3f945647ffb95b5839c07038d64f9811adf17308b9121d8a2b87b6a22a80a39 AS runtime

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
COPY --from=build /go/bin/example-health /example-health
