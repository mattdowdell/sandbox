package step

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/tests/utils/authnv1client"
	"github.com/mattdowdell/sandbox/tests/utils/examplev1client"
	"github.com/mattdowdell/sandbox/tests/utils/input"
)

// ...
func PrintableASCIIChars(ctx context.Context, length int) (context.Context, error) {
	name, err := randomString(printableASCII(), length)
	if err != nil {
		return ctx, err
	}

	return input.NameIntoContext(ctx, name), nil
}

// ...
func PrintableNonASCIIChars(ctx context.Context, length int) (context.Context, error) {
	name, err := randomString(printableNonASCII(), length)
	if err != nil {
		return ctx, err
	}

	return input.NameIntoContext(ctx, name), nil
}

func ExistingResourceName(ctx context.Context) (context.Context, error) {
	//nolint:mnd // arbitrary length chosen
	name, err := randomString(printableASCII(), 20 /*length*/)
	if err != nil {
		return ctx, err
	}

	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	_, cleanup, err := client.CreateResource(ctx, name)
	if err != nil {
		return ctx, err // TODO: wrap error?
	}

	ctx = examplev1client.AppendCleanup(ctx, cleanup)
	ctx = input.NameIntoContext(ctx, name)

	return ctx, nil
}

// EmptyID inserts an empty string as an ID into the context.
func EmptyID(ctx context.Context) context.Context {
	return input.IDIntoContext(ctx, "")
}

// NilUUID inserts a nil UUID as an ID into the context.
func NilUUID(ctx context.Context) context.Context {
	return input.IDIntoContext(ctx, uuid.Nil.String())
}

// InvalidUUID inserts an invalid UUID as an ID into the context.
func InvalidUUID(ctx context.Context) context.Context {
	return input.IDIntoContext(ctx, "invalid")
}

// ExistingID creates a new Resource and adds its ID into the context.
func ExistingID(ctx context.Context) (context.Context, error) {
	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	//nolint:mnd // arbitrary length chosen
	name, err := randomString(printableASCII(), 20 /*length*/)
	if err != nil {
		return ctx, err
	}

	resource, cleanup, err := client.CreateResource(ctx, name)
	if err != nil {
		return ctx, err
	}

	ctx = examplev1client.AppendCleanup(ctx, cleanup)
	ctx = input.IDIntoContext(ctx, resource.GetId())
	ctx = input.NameIntoContext(ctx, name)

	return ctx, err
}

// ExistingResources creates the specified number of Resources and adds them to the context.
func ExistingResources(ctx context.Context, count int) (context.Context, error) {
	client, err := examplev1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	cleanups := make([]examplev1client.Cleanup, 0, count)

	for range count {
		//nolint:mnd // arbitrary length chosen
		name, err := randomString(printableASCII(), 20 /*length*/)
		if err != nil {
			return ctx, err
		}

		_, cleanup, err := client.CreateResource(ctx, name)
		if err != nil {
			return ctx, err
		}

		cleanups = append(cleanups, cleanup)
	}

	return examplev1client.AppendCleanup(ctx, cleanups...), nil
}

// ...
func Limit(ctx context.Context, value int32) context.Context {
	return input.LimitIntoContext(ctx, value)
}

// ...
func InvalidAuthentication(ctx context.Context) context.Context {
	return input.AuthnIntoContext(ctx, "Basic invalid")
}

// ...
func NoAuthentication(ctx context.Context) context.Context {
	return input.AuthnIntoContext(ctx, "")
}

// ...
func ValidAuthentication(ctx context.Context) (context.Context, error) {
	client, err := authnv1client.FromContext(ctx)
	if err != nil {
		return ctx, err
	}

	token, err := client.Login(ctx, "username", "password")
	if err != nil {
		return ctx, err
	}

	return input.AuthnIntoContext(ctx, "Bearer "+token), nil
}

// ...
func NonExistingConfigKey(ctx context.Context) (context.Context, error) {
	return input.ConfigKeyIntoContext(ctx, "does_not_exist"), nil
}

// ...
func EmptyConfigKey(ctx context.Context) (context.Context, error) {
	return input.ConfigKeyIntoContext(ctx, ""), nil
}

// ...
func ExistingConfigKey(ctx context.Context, key string) (context.Context, error) {
	return input.ConfigKeyIntoContext(ctx, key), nil
}

// printableASCII returns a set of printable ASCII characters for use with RandomString.
func printableASCII() string {
	var output []rune

	for c := ' '; c <= '~'; c++ {
		output = append(output, c)
	}

	return string(output)
}

// printableNonASCII returns a set of printable non-ASCII characters for use with RandomString.
//
// Specifically, it returns the characters code from Latin Extended-A and Latin Extended-B.
func printableNonASCII() string {
	var output []rune

	for c := 'Ā'; c <= 'ɏ'; c++ {
		output = append(output, c)
	}

	return string(output)
}

// randomString generates a string of the given length using the set of given characters.
func randomString(chars string, length int) (string, error) {
	in := []rune(chars)
	out := make([]rune, 0, length)

	l := big.NewInt(int64(len(in)))

	for range length {
		i, err := rand.Int(rand.Reader, l)
		if err != nil {
			return "", err
		}

		out = append(out, in[i.Int64()])
	}

	return string(out), nil
}
