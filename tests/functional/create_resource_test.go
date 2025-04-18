package functional_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/cucumber/godog"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/gen/example/v1/examplev1connect"
	"github.com/mattdowdell/sandbox/tests/utils"
)

func InitCreateResource(sc *godog.ScenarioContext) {
	c := NewCreateResource()

	sc.Given(`^a name of (\d+) printable ASCII characters$`, c.aNameOfPrintableASCIICharacters)
	sc.Given(`^a name of (\d+) printable non-ASCII characters$`, c.aNameOfPrintableNonASCIICharacters)
	sc.Given(`^an existing resource name$`, c.anExistingResourceName)
	sc.When(`^I create a Resource$`, c.iCreateAResource)
	sc.Then(`^I should receive AlreadyExists with message: (.+)$`, c.iShouldReceiveAlreadyExists)
	sc.Then(`^I should receive InvalidArgument with message: (.+)$`, c.iShouldReceiveInvalidArgument)
	sc.Then(`^I should receive the Resource$`, c.iShouldReceiveTheResource)

	sc.After(c.afterScenarioHook)
}

type CreateResource struct {
	client    examplev1connect.ExampleServiceClient
	validator protovalidate.Validator
	name      string
	resource  *examplev1.Resource
	err       error
}

func NewCreateResource() *CreateResource {
	client := examplev1connect.NewExampleServiceClient(
		http.DefaultClient,
		"http://localhost:5000", // TODO: pull from config
	)

	validator, err := protovalidate.New()
	if err != nil {
		panic(err) // TODO: handle this better
	}

	return &CreateResource{
		client:    client,
		validator: validator,
	}
}

func (c *CreateResource) afterScenarioHook(
	ctx context.Context,
	_ *godog.Scenario,
	err error,
) (context.Context, error) {
	id := c.resource.GetId()
	if id == "" {
		return ctx, err
	}

	if _, err2 := c.client.DeleteResource(
		ctx,
		connect.NewRequest(&examplev1.DeleteResourceRequest{
			Id: id,
		}),
	); err2 != nil {
		return ctx, errors.Join(err, err2)
	}

	return ctx, err
}

func (c *CreateResource) aNameOfPrintableASCIICharacters(length int) error {
	name, err := utils.RandomString(utils.PrintableASCII(), length)
	if err != nil {
		return err
	}

	c.name = name
	return nil
}

func (c *CreateResource) aNameOfPrintableNonASCIICharacters(length int) error {
	name, err := utils.RandomString(utils.PrintableNonASCII(), length)
	if err != nil {
		return err
	}

	c.name = name

	return nil
}

func (c *CreateResource) anExistingResourceName() error {
	name, err := utils.RandomString(utils.PrintableASCII(), 20)
	if err != nil {
		return err
	}

	c.name = name

	resp, err := c.client.CreateResource(
		context.Background(),
		connect.NewRequest(&examplev1.CreateResourceRequest{
			Resource: &examplev1.ResourceCreate{
				Name: c.name,
			},
		}),
	)
	if err != nil {
		return err // TODO: wrap
	}

	if err := c.validator.Validate(resp.Msg); err != nil {
		return err // TODO: wrap
	}

	return nil
}

func (c *CreateResource) iCreateAResource() error {
	resp, err := c.client.CreateResource(
		context.Background(),
		connect.NewRequest(&examplev1.CreateResourceRequest{
			Resource: &examplev1.ResourceCreate{
				Name: c.name,
			},
		}),
	)
	if err != nil {
		c.err = err
		return nil //nolint:nilerr // to support tests that expect an error
	}

	if err := c.validator.Validate(resp.Msg); err != nil {
		return err
	}

	c.resource = resp.Msg.GetResource()
	return nil
}

func (c *CreateResource) iShouldReceiveAlreadyExists(msg string) error {
	if err := utils.CheckConnectCode(c.err, connect.CodeAlreadyExists); err != nil {
		return err
	}

	return utils.CheckConnectMsg(c.err, msg)
}

func (c *CreateResource) iShouldReceiveInvalidArgument(msg string) error {
	if err := utils.CheckConnectCode(c.err, connect.CodeInvalidArgument); err != nil {
		return err
	}

	return utils.CheckConnectMsg(c.err, msg)
}

func (c *CreateResource) iShouldReceiveTheResource() error {
	var errs []error

	if n := c.resource.GetName(); n != c.name {
		errs = append(errs, fmt.Errorf("unexpected name: want: %q, have: %q", c.name, n))
	}

	return errors.Join(errs...)
}
