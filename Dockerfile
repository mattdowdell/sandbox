# https://docs.docker.com/reference/dockerfile/
#
# TODO: document use of mirror.gcr.io
# https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images

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

FROM base AS dev

RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends tini; \
    rm -rf /var/lib/apt/lists/*;

RUN set -eux; \
    curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | \
    bash -s -- --to /usr/local/bin;

ARG TRIVY_VERSION=v0.59.1
RUN set -eux; \
    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | \
    sh -s -- -b /usr/local/bin ${TRIVY_VERSION};

ARG USER_NAME=dev
ARG USER_COMMENT=Dev
ARG USER_UID=1000
ARG USER_GID=1000
RUN set -eux; \
    groupadd --gid ${USER_GID} ${USER_NAME}; \
    useradd --comment ${USER_COMMENT} --create-home --gid ${USER_GID} --uid ${USER_UID} ${USER_NAME}

ENV GOPATH=''
USER ${USER_NAME}
RUN set -eux; \
    echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> /home/dev/.bashrc

ARG GOLANGCI_LINT_VERSION=v1.64.5
RUN set -eux; \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION};

WORKDIR /home/dev/ws

CMD ["sleep", "inf"]
