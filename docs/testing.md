# Testing

## Unit Testing

_TODO: document unit tests philosophy, e.g. no IO, etc._

### Writing Tests

_TODO: document unit testing style._

### Execution

The unit tests can be executed by running:

```sh
just unit
```

By default, a timeout of 30 seconds is applied. This can be adjusted with an additional argument:

```sh
# change to a 10 minute timeout
just unit 10m
```

### Coverage

_TODO: document how to view and interpret coverage results_

## Functional Testing

Functional tests exist to check that the service's RPC API is working correctly. This includes
validating success responses and client error responses. It does not include validation of server
error responses as they are not currently possible to induce.

Functional testing follows [Behaviour-Driven Development] (BDD), with test scenarios defined using
[Gherkin]. Scenarios are grouped into features, with each RPC method having a separate feature. The
current set of features can be found in `./tests/functional/features`.

Gherkin was chosen because it's structured language is easy for both developers and non-developers
to understand. It allows tests to be written alongside the design, where multiple stakeholders can
review the test plan for both correctness and completeness.

Furthermore, the use of Gherkin enables steps within a scenario to be re-used with minimal effort.
Once a step has been implemented, it can often be added to other scenarios across multiple features
without touching the step implementation. This re-use reduces the complexity and maintenance burden
of the code needed to test the service.

[Behaviour-Driven Development]: https://cucumber.io/docs/bdd/
[Gherkin]: https://cucumber.io/docs/gherkin/reference

### Writing Tests

_TODO: document how to add new tests._

### Execution

The functional tests can be executed by running:

```sh
just functional
```

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
