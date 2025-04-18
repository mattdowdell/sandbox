package utils

import (
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
)

// ...
func CheckConnectCode(err error, want connect.Code) error {
	cast, err := castConnectError(err)
	if err != nil {
		return err
	}

	if cast.Code() == want {
		return nil
	}

	return fmt.Errorf(
		"unexpected code: want: %d (%s), have: %d (%v)",
		want,
		want.String(),
		cast.Code(),
		cast,
	)
}

// ...
func CheckConnectMsg(err error, want string) error {
	cast, err := castConnectError(err)
	if err != nil {
		return err
	}

	// TODO: figure out how to represent newlines in feature files
	have := strings.ReplaceAll(cast.Message(), "\n", "")

	if have != want {
		return fmt.Errorf("unexpected msg: want: %q, have: %q", want, have)
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
