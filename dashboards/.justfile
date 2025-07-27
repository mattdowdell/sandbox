# https://just.systems/man/en/

import '../.tools.just'

[private]
default:
    @just --list

# TODO
build: build-connectrpc

# TODO
build-connectrpc:
    {{ jsonnet }} -J vendor ./connectrpc/dashboard.jsonnet > ../config/dashboards/connectrpc_generated.json
