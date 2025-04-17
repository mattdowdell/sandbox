# sandbox

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

[ConnectRPC]: #todo
[protovalidate]: #todo
[Buf]: #todo
[Jet]: #todo
[OpenTelemetry]: #todo
[Koanf]: #todo

## Sanity tests

```sh
# ready: success
echo '{}' | grpc-client-cli -a localhost:5000 -s Health -m Check

# create: invalid argument
echo '{}' | grpc-client-cli -a localhost:5000 -s ExampleService -m CreateResource

# create: success
echo '{"resource":{"name":"example"}}' | \
	grpc-client-cli -a localhost:5000 -s ExampleService -m CreateResource

# create: success (bulk)
for i in {1..100}; do
	echo "{\"resource\":{\"name\":\"example-$i\"}}" | \
		grpc-client-cli -a localhost:5000 -s ExampleService -m CreateResource
done

# list: success
echo '{}' | grpc-client-cli -a localhost:5000 -s ExampleService -m ListResources

# get: not found
echo '{"id":"00000000-0000-0000-0000-000000000000"}' | \
	grpc-client-cli -a localhost:5000 -s ExampleService -m GetResource

# update: not found
echo '{"resource":{"id":"00000000-0000-0000-0000-000000000000","name":"example"}}' | \
        grpc-client-cli -a localhost:5000 -s ExampleService -m UpdateResource

# update
echo '{"resource":{"id":"01961ee2-68d4-7708-9c97-9297e402ab59","name":"example"}}' | \
	grpc-client-cli -a localhost:5000 -s ExampleService -m UpdateResource

# delete: not found
echo '{"id":"00000000-0000-0000-0000-000000000000"}' | \
        grpc-client-cli -a localhost:5000 -s ExampleService -m DeleteResource
```
