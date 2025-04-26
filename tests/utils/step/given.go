package step

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"

	"github.com/mattdowdell/sandbox/tests/utils/examplev1client"
	"github.com/mattdowdell/sandbox/tests/utils/input"
)

// ...
func PrintableASCIIChars(ctx context.Context, length int) (context.Context, error) {
	name, err := randomString(printableASCII(), length)
	if err != nil {
		return ctx, err
	}

	return input.AddNameToContext(ctx, name), nil
}

// ...
func PrintableNonASCIIChars(ctx context.Context, length int) (context.Context, error) {
	name, err := randomString(printableNonASCII(), length)
	if err != nil {
		return ctx, err
	}

	return input.AddNameToContext(ctx, name), nil
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
	ctx = input.AddNameToContext(ctx, name)

	return ctx, nil
}

// ...
func NilUUID(ctx context.Context) context.Context {
	return input.AddIDToContext(ctx, uuid.Nil.String())
}

// ...
func InvalidUUID(ctx context.Context) context.Context {
	return input.AddIDToContext(ctx, "invalid")
}

// ...
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
	ctx = input.AddIDToContext(ctx, resource.GetId())
	ctx = input.AddNameToContext(ctx, name)

	return ctx, err
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
	input := []rune(chars)
	output := make([]rune, 0, length)

	l := big.NewInt(int64(len(input)))

	for range length {
		i, err := rand.Int(rand.Reader, l)
		if err != nil {
			return "", err
		}

		output = append(output, input[i.Int64()])
	}

	return string(output), nil
}
