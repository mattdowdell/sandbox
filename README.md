# sandbox

[![Documentation](https://godoc.org/github.com/mattdowdell/sandbox?status.svg)](http://godoc.org/github.com/mattdowdell/sandbox)
[![License](https://img.shields.io/github/license/mattdowdell/sandbox)](https://github.com/mattdowdell/sandbox?tab=MIT-1-ov-file#readme)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/mattdowdell/sandbox?label=version)](https://github.com/mattdowdell/sandbox/tags)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mattdowdell/sandbox/ci.yml)](https://github.com/mattdowdell/sandbox/actions/workflows/ci.yml)

A toy Go microservice intended for use as reference material.

## Features

- Uses a directory structure (mostly) adherent to clean architecture.
- API provided by [ConnectRPC], validation using [protovalidate], [Buf] generated, reflection enabled.
- Database provided by PostgreSQL, UUIDv7 for primary keys, [Jet] for SQL query building.
- [OpenTelemetry] Metrics and Tracing, standardised metrics and attributes whenever possible.
- Logging using `log/slog`.
- Runtime configuration via [Koanf].
- Grafana (with version-controlled dashboards), VictoriaMetrics and Jaeger for observability.
- CI using GitHub actions, with Zizmor for linting.
- Vulnerability and License scanning using Trivy.
- Secret scanning using Gitleaks.
- Dependency updates from Dependabot.
- Packaged using Docker containers, developed with Docker Compose.

[ConnectRPC]: https://github.com/connectrpc/connect-go
[protovalidate]: https://github.com/bufbuild/protovalidate
[Buf]: https://github.com/bufbuild/buf
[Jet]: https://pkg.go.dev/github.com/go-jet/jet/v2
[OpenTelemetry]: https://opentelemetry.io
[Koanf]: https://pkg.go.dev/github.com/knadh/koanf/v2

## Sanity tests

```sh
# ready: success
echo '{}' | grpc-client-cli -a localhost:5000 -s Health -m Check

# login
echo '{"id":"example","secret":"example"}' | grpc-client-cli -a localhost:5000 -s AuthnService -m Login

token=`echo '{"id":"example","secret":"example"}' | \
  grpc-client-cli -a localhost:5000 -s AuthnService -m Login | \
  jq -r '.access_token'`

# authenticate
echo "{\"token\":\"${token>}\"}" | grpc-client-cli -a localhost:5000 -s AuthnService -m Authenticate

# create: success
echo '{"resource":{"name":"example"}}' | \
	grpc-client-cli \
		-a localhost:5000 \
		-s ExampleService \
		-m CreateResource \
		-H "Authorization: Bearer $token"

# create: success (bulk)
for i in {1..100}; do
	echo "{\"resource\":{\"name\":\"example-$i\"}}" | \
		grpc-client-cli \
		-a localhost:5000 \
		-s ExampleService \
		-m CreateResource \
		-H "Authorization: Bearer $token"
done

# list: success
echo '{"limit": 100}' | \
	grpc-client-cli \
		-a localhost:5000 \
		-s ExampleService \
		-m ListResources \
		-H "Authorization: Bearer $token"

```
