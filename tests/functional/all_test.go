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

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	o := opts
	o.TestingT = t

	status := godog.TestSuite{
		Name:                 "create_resource",
		Options:              &o,
		// TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
	}.Run()

	// if status == 2 {
	// 	t.SkipNow()
	// }

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}
