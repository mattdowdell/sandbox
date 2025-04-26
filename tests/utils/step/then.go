package step

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"

	"github.com/mattdowdell/sandbox/tests/utils/input"
	"github.com/mattdowdell/sandbox/tests/utils/output"
)

// ...
func FailWithCodeAndMsg(ctx context.Context, code, msg string) error {
	have, err := output.ErrFromContext(ctx)
	if err != nil {
		return err
	}

	return checkConnectErr(have, code, msg)
}

// ...
func CheckResource(ctx context.Context) error {
	name, err := input.NameFromContext(ctx)
	if err != nil {
		return err
	}

	resource, err := output.ResourceFromContext(ctx)
	if err != nil {
		return err
	}

	var errs []error

	if want, err := input.IDFromContext(ctx); err == nil {
		if have := resource.GetId(); want != have {
			errs = append(errs, fmt.Errorf("unexpected id: want: %q, have: %q", want, have))
		}
	}

	if have := resource.GetName(); have != name {
		errs = append(errs, fmt.Errorf("unexpected name: want: %q, have: %q", name, have))
	}

	return errors.Join(errs...)
}

// ...
func CheckSuccess(ctx context.Context) error {
	return output.EmptyFromContext(ctx)
}

// ...
func checkConnectErr(err error, code, msg string) error {
	cast, err := castConnectError(err)
	if err != nil {
		return err
	}

	var wantCode connect.Code
	if err := wantCode.UnmarshalText([]byte(code)); err != nil {
		return err
	}

	if cast.Code() != wantCode {
		return fmt.Errorf(
			"unexpected code: want: %d (%s), have: %d (%v)",
			wantCode,
			wantCode.String(),
			cast.Code(),
			cast,
		)
	}

	// TODO: figure out how to represent newlines in feature files
	wantMsg := strings.ReplaceAll(cast.Message(), "\n", "")

	if wantMsg != msg {
		return fmt.Errorf("unexpected msg: want: %q, have: %q", wantMsg, msg)
	}

	return nil
}

// ...
func castConnectError(err error) (*connect.Error, error) {
	var cast *connect.Error
	if errors.As(err, &cast) {
		return cast, nil
	}

	return nil, fmt.Errorf(
		"unexpected error type: want: %T, have: %T (%v)",
		&connect.Error{},
		err,
		err,
	)
}
