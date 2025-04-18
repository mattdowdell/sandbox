package functional_test

import (
	"flag"
	"testing"

	"github.com/cucumber/godog"
)

/*
TODO: The below suites are an outline of the success and error scenarios that can reasonably be
triggered.

As this is a low-risk project, investigate using BDD with godog to see how that compares to
more Go-idiomatic testing approaches.

Key requirements:
- Tests can be run from a single binary and entrypoint, enabling of istio sidecar shutdown.
- Useful errors to debug with.
- Selection of individual tests (maybe?)
*/

// func Test_All(t *testing.T) {
// 	suites := []suite.TestingSuite{
// 		NewCreateResource(),
// 		NewGetResource(),
// 		NewUpdateResource(),
// 		NewDeleteResource(),
// 		NewListResources(),
// 		NewListAuditEvents(),
// 		NewWatchAuditEvents(),
// 	}

// 	for _, s := range suites {
// 		suite.Run(t, s)
// 	}
// }

var opts = godog.Options{}

//nolint:gochecknoinits // only way to configure flags
func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestCreateResource(t *testing.T) {
	o := opts
	o.TestingT = t

	suite := godog.TestSuite{
		Name:                "create_resource",
		ScenarioInitializer: InitCreateResource,
		Options:             &opts,
	}

	if status := suite.Run(); status != 0 {
		t.Fatalf("non-zero status code received: %d", status)
	}
}
