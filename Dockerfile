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

# -------------
# Devenv target
# -------------

FROM mirror.gcr.io/golangci/golangci-lint:v1.64.5@sha256:9faef4dda4304c4790a14c5b8c8cd8c2715a8cb754e13f61d8ceaa358f5a454a AS golangci-lint

FROM base AS dev

RUN set -eux; \
    curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | \
    bash -s -- --to /usr/local/bin;

ARG TRIVY_VERSION=v0.59.1
RUN set -eux; \
    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | \
    sh -s -- -b /usr/local/bin ${TRIVY_VERSION};

RUN set -eux; \
    useradd --comment Dev --create-home --user-group dev;

ENV GOPATH=''
USER dev
RUN set -eux; \
    echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> /home/dev/.bashrc; \
    mkdir -p /home/dev/ws; \
    git config --global --add safe.directory /home/dev/ws;

COPY --from=golangci-lint /usr/bin/golangci-lint /usr/bin/golangci-lint

WORKDIR /home/dev/ws

CMD ["sleep", "inf"]
