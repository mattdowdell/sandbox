# TODO: link docs
#
# TODO: document use of mirror.gcr.io
# https://cloud.google.com/artifact-registry/docs/pull-cached-dockerhub-images

# ------------
# Build target
# ------------

FROM mirror.gcr.io/golang:1.23-bookworm@sha256:2e838582004fab0931693a3a84743ceccfbfeeafa8187e87291a1afea457ff7a AS build

WORKDIR /go/src

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-buildid= -s -w" -o /go/bin/ ./cmd/...;

# --------------
# Runtime target
# --------------

FROM gcr.io/distroless/static-debian12:nonroot@sha256:6cd937e9155bdfd805d1b94e037f9d6a899603306030936a3b11680af0c2ed58 AS runtime

ARG SERVICE
COPY --from=build /go/bin/${SERVICE} /${SERVICE}
