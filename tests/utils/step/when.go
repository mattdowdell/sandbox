package step

import (
	"context"

	"github.com/mattdowdell/sandbox/tests/utils/examplev1client"
	"github.com/mattdowdell/sandbox/tests/utils/input"
	"github.com/mattdowdell/sandbox/tests/utils/output"
)

// ...
func CreateResource(ctx context.Context) (context.Context, error) {
	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	name, err := input.NameFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	resource, cleanup, err := client.CreateResource(ctx, name)
	if err != nil {
		return output.AddErrToContext(ctx, err), nil
	}

	ctx = examplev1client.AppendCleanup(ctx, cleanup)
	ctx = output.AddResourceToContext(ctx, resource)

	return ctx, nil
}

// ...
func GetResource(ctx context.Context) (context.Context, error) {
	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	id, err := input.IDFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	resource, err := client.GetResource(ctx, id)
	if err != nil {
		return output.AddErrToContext(ctx, err), nil
	}

	return output.AddResourceToContext(ctx, resource), nil
}

// ...
func DeleteResource(ctx context.Context) (context.Context, error) {
	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	id, err := input.IDFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	if err := client.DeleteResource(ctx, id); err != nil {
		return output.AddErrToContext(ctx, err), nil
	}

	return output.AddEmptyToContext(ctx), nil
}
