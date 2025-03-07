package functional_test

import (
	"net/http"

	"github.com/cucumber/godog"
	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/gen/example/v1/examplev1connect"
	"github.com/mattdowdell/sandbox/tests/utils"
)

type CreateResource struct {
	client examplev1connect.ExampleServiceClient
	name string
	resource *examplev1.Resource
	err error
}

func NewCreateResource() *CreateResource {
	client := examplev1connect.NewExampleServiceClient(
		http.DefaultClient,
		"http://localhost:5000", // TODO: pull from config
	)

	return &CreateResource{
		client: client,
	}
}

func (c *CreateResource) aNameOfPrintableASCIICharacters(length int) error {
	name, err := utils.RandomString(utils.PrintableASCII(), length)
	if err != nil {
		return err
	}

	c.name = name
	return nil
}

func (c *CreateResource) aNameOfPrintableNonASCIICharacters(arg1 int) error {
	return godog.ErrPending
}

func (c *CreateResource) anExistingResourceName() error {
	return godog.ErrPending
}

func (c *CreateResource) iCreateAResource() error {
	return nil
}

func (c *CreateResource) iShouldReceiveAlreadyExists() error {
	return godog.ErrPending
}

func (c *CreateResource) iShouldReceiveInvalidArgument() error {
	return godog.ErrPending
}

func (c *CreateResource) iShouldReceiveTheResource() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	c := NewCreateResource()

	ctx.Step(`^a name of (\d+) printable ASCII characters$`, c.aNameOfPrintableASCIICharacters)
	ctx.Step(`^a name of (\d+) printable non-ASCII characters$`, c.aNameOfPrintableNonASCIICharacters)
	ctx.Step(`^an existing resource name$`, c.anExistingResourceName)
	ctx.Step(`^I create a Resource$`, c.iCreateAResource)
	ctx.Step(`^I should receive AlreadyExists$`, c.iShouldReceiveAlreadyExists)
	ctx.Step(`^I should receive InvalidArgument$`, c.iShouldReceiveInvalidArgument)
	ctx.Step(`^I should receive the Resource$`, c.iShouldReceiveTheResource)
}
