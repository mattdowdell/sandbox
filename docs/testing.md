# Testing

## Unit testing

### Coverage

## Functional Testing

### Coverage

Functional test coverage is not enabled by default. However, it can be collected by enabling
coverage profiling in the Development Environment's rpc container build, running the tests and then
stopping the development environment.

```sh
# delete any stale artifacts
just functional-cover-clean

# enable coverage profiling
echo GO_BUILD_ARGS=-cover > .env

# rebuild the rpc container
just dev-restart

# execute the functional tests
# TODO

# shutdown the rpc container to output coverage
just dev-down
```

This will produce a number of files in the `.covdata/` directory which can be coverted to a HTML
report:

```sh
# output the overall coverage and generate `functional.html`
just functional-cover
```
