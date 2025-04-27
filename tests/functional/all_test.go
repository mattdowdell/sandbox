package functional_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/cucumber/godog"

	"github.com/mattdowdell/sandbox/tests/utils/examplev1client"
	"github.com/mattdowdell/sandbox/tests/utils/step"
)

var opts godog.Options

//nolint:gochecknoinits // recommended way to support runtime customisation
func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func Test_ExampleService(t *testing.T) {
	client := examplev1client.New("http://localhost:5000")

	o := opts
	o.TestingT = t
	o.DefaultContext = examplev1client.AddToContext(t.Context(), client)

	suite := godog.TestSuite{
		Name:                "example_service",
		Options:             &o,
		ScenarioInitializer: InitializeScenario,
	}

	switch status := suite.Run(); status {
	case 0:
		// success
	case 2:
		t.SkipNow()
	default:
		t.Fatalf("zero status code expected, %d received", status)
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, scen *godog.Scenario) (context.Context, error) {
		fmt.Println("name:", scen.Name)
		return ctx, nil
	})

	sc.Given(`^a name of (\d+) printable ASCII characters$`, step.PrintableASCIIChars)
	sc.Given(`^a name of (\d+) printable non-ASCII characters$`, step.PrintableNonASCIIChars)
	sc.Given(`^an existing resource name$`, step.ExistingResourceName)
	sc.Given(`^a non-existent Resource ID$`, step.NilUUID)
	sc.Given(`^an invalid Resource ID$`, step.InvalidUUID)
	sc.Given(`^an existing Resource ID$`, step.ExistingID)

	sc.When(`^I create a Resource$`, step.CreateResource)
	sc.When(`^I get the Resource$`, step.GetResource)
	sc.When(`^I update the Resource$`, step.UpdateResource)
	sc.When(`^I delete the Resource$`, step.DeleteResource)

	sc.Then(`^I should fail with code=(\w+), msg=(.+)$$`, step.FailWithCodeAndMsg)
	sc.Then(`^I should receive the Resource$`, step.CheckResource)
	sc.Then(`^I should succeed$`, step.CheckSuccess)

	sc.After(func(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {
		if err2 := examplev1client.RunCleanups(ctx); err2 != nil {
			return ctx, errors.Join(err, err2)
		}

		return ctx, nil
	})
}
