# https://docs.docker.com/reference/dockerfile/
#
# mirror.gcr.io caches of popular docker hub images, but does not add rate limiting.
# See https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images.

# -----------
# Base target
# -----------

FROM mirror.gcr.io/golang:1.23-bookworm@sha256:6260304a09fb81a1983db97c9e6bfc1779ebce33d39581979a511b3c7991f076 AS base

# ------------
# Build target
# ------------

FROM base AS build

WORKDIR /go/src

RUN --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -trimpath -ldflags="-buildid= -s -w" -o /go/bin/ ./cmd/...;

# --------------
# Runtime target
# --------------

FROM gcr.io/distroless/static-debian12:nonroot@sha256:6ec5aa99dc335666e79dc64e4a6c8b89c33a543a1967f20d360922a80dd21f02 AS runtime

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
COPY --from=build /go/bin/example-health /example-health
