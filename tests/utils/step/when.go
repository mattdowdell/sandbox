package step

import (
	"context"

	"github.com/mattdowdell/sandbox/tests/utils/configv1client"
	"github.com/mattdowdell/sandbox/tests/utils/examplev1client"
	"github.com/mattdowdell/sandbox/tests/utils/input"
	"github.com/mattdowdell/sandbox/tests/utils/output"
)

// CreateResource creates a Resource using a name input taken from the given context.
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
		return output.ErrIntoContext(ctx, err), nil
	}

	ctx = examplev1client.AppendCleanup(ctx, cleanup)
	ctx = output.ResourceIntoContext(ctx, resource)

	return ctx, nil
}

// GetResource gets a single Resource using the ID input taken from the given context.
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
		return output.ErrIntoContext(ctx, err), nil
	}

	return output.ResourceIntoContext(ctx, resource), nil
}

// ListResources gets multiple Resources using parameters taken from the given context.
func ListResources(ctx context.Context) (context.Context, error) {
	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	limit, err := input.LimitFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	resources, next, err := client.ListResources(ctx, limit)
	if err != nil {
		return output.ErrIntoContext(ctx, err), nil
	}

	ctx = output.ResourcesIntoContext(ctx, resources)
	ctx = output.NextIntoContext(ctx, next)

	return ctx, nil
}

// UpdateResource updates a Resource using the ID and name taken from the given context.
func UpdateResource(ctx context.Context) (context.Context, error) {
	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	id, err := input.IDFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	name, err := input.NameFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	resource, err := client.UpdateResource(ctx, id, name)
	if err != nil {
		return output.ErrIntoContext(ctx, err), nil
	}

	return output.ResourceIntoContext(ctx, resource), nil
}

// DeleteResource deletes a Resource using the ID taken from the given context.
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
		return output.ErrIntoContext(ctx, err), nil
	}

	return output.EmptyIntoContext(ctx), nil
}

// GetConfigValue gets a single configuration value using the key taken from the given context.
func GetConfigValue(ctx context.Context) (context.Context, error) {
	client, err := configv1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	key, err := input.ConfigKeyFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	value, err := client.GetConfigValue(ctx, key)
	if err != nil {
		return output.ErrIntoContext(ctx, err), nil
	}

	return output.ConfigValueIntoContext(ctx, value), nil
}
