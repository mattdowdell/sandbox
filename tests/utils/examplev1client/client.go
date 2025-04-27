package examplev1client

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/gen/example/v1"
	"github.com/mattdowdell/sandbox/gen/example/v1/examplev1connect"
)

// ...
type Cleanup = func(context.Context) error

// ...
type Client struct {
	inner examplev1connect.ExampleServiceClient
}

// ...
func New(baseURL string) *Client {
	inner := examplev1connect.NewExampleServiceClient(
		http.DefaultClient,
		baseURL,
		connect.WithInterceptors(
			connect.UnaryInterceptorFunc(ValidateUnaryInterceptor),
			connect.UnaryInterceptorFunc(ScenarioUnaryInterceptor),
		),
	)

	return &Client{
		inner: inner,
	}
}

// ...
func (c *Client) CreateResource(ctx context.Context, name string) (*examplev1.Resource, Cleanup, error) {
	resp, err := c.inner.CreateResource(
		ctx,
		connect.NewRequest(&examplev1.CreateResourceRequest{
			Resource: &examplev1.ResourceCreate{
				Name: name,
			},
		}),
	)
	if err != nil {
		return nil, nil, err
	}

	resource := resp.Msg.GetResource()
	return resource, c.cleanup(resource.GetId()), nil
}

// ...
func (c *Client) GetResource(ctx context.Context, id string) (*examplev1.Resource, error) {
	resp, err := c.inner.GetResource(
		ctx,
		connect.NewRequest(&examplev1.GetResourceRequest{
			Id: id,
		}),
	)
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetResource(), nil
}

// ...
func (c *Client) UpdateResource(ctx context.Context, id, name string) (*examplev1.Resource, error) {
	resp, err := c.inner.UpdateResource(
		ctx,
		connect.NewRequest(&examplev1.UpdateResourceRequest{
			Resource: &examplev1.ResourceUpdate{
				Id:   id,
				Name: name,
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetResource(), nil
}

// ...
func (c *Client) DeleteResource(ctx context.Context, id string) error {
	_, err := c.inner.DeleteResource(
		ctx,
		connect.NewRequest(&examplev1.DeleteResourceRequest{
			Id: id,
		}),
	)
	return err
}

func (c *Client) cleanup(id string) Cleanup {
	return func(ctx context.Context) error {
		err := c.DeleteResource(ctx, id)
		if err == nil {
			return nil
		}

		var cast *connect.Error
		if errors.As(err, &cast) && cast.Code() == connect.CodeNotFound {
			return nil
		}

		return err
	}
}
