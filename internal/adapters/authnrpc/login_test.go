package authnrpc_test

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/assert"

	authnv1 "github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/authnrpc"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockauthnrpc"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

const (
	testID        = "id"
	testSecret    = "secret"
	testToken     = "an.example.jwt" //nolint:gosec // not a real secret
	testTokenType = "Bearer"
	testExpiresIn = uint32(7200)
)

func Test_Handler_Login_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	issuer := mockauthnrpc.NewIssuer(t)
	issuer.EXPECT().Issue(testID, testExpiresIn).Return(testToken, nil).Once()

	handler := authnrpc.New(
		issuer,
		mockauthnrpc.NewParser(t),
	)

	req := connect.NewRequest(&authnv1.LoginRequest{
		Id:     testID,
		Secret: testSecret,
	})

	// act
	resp, err := handler.Login(ctx, req)

	// assert
	expected := connect.NewResponse(&authnv1.LoginResponse{
		AccessToken: testToken,
		TokenType:   testTokenType,
		ExpiresIn:   testExpiresIn,
	})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Handler_Login_Error(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	issuer := mockauthnrpc.NewIssuer(t)
	issuer.EXPECT().Issue(testID, testExpiresIn).Return("", errors.New("example")).Once()

	handler := authnrpc.New(
		issuer,
		mockauthnrpc.NewParser(t),
	)

	req := connect.NewRequest(&authnv1.LoginRequest{
		Id:     testID,
		Secret: testSecret,
	})

	// act
	resp, err := handler.Login(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}
