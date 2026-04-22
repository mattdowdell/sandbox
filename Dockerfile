# https://docs.docker.com/reference/dockerfile/
#
# mirror.gcr.io caches of popular docker hub images, but does not add rate limiting.
# See https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images.

# ------------
# Build target
# ------------

FROM --platform=$BUILDPLATFORM mirror.gcr.io/golang:1.26-trixie@sha256:cd8540d626ab35272a8f5ef829c5e5f189ce1dfbf467c123320c191c69c23245 AS build

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

FROM gcr.io/distroless/static-debian13:debug-nonroot@sha256:fa2ca4ab94779411a1d20c2f5205ae550f1a366267b3ae6dfe9c40faae1bcb1e AS debug

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
