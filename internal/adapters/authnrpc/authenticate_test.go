package authnrpc_test

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"

	authnv1 "github.com/mattdowdell/sandbox/gen/authn/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/authnrpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_Handler_Authenticate_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	claims := authnrpc.NewMockClaims(t)
	claims.EXPECT().GetSubject().Return(testID, nil).Once()

	parser := authnrpc.NewMockParser(t)
	parser.EXPECT().Parse(testToken).Return(claims, nil).Once()

	handler := authnrpc.New(
		authnrpc.NewMockIssuer(t),
		parser,
	)

	req := connect.NewRequest(&authnv1.AuthenticateRequest{
		Token: testToken,
	})

	// act
	resp, err := handler.Authenticate(ctx, req)

	// assert
	expected := connect.NewResponse(&authnv1.AuthenticateResponse{
		Subject: testID,
	})

	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
}

func Test_Handler_Authenticate_ParseError(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	parser := authnrpc.NewMockParser(t)
	parser.EXPECT().Parse(testToken).Return(nil, errors.New("example")).Once()

	handler := authnrpc.New(
		authnrpc.NewMockIssuer(t),
		parser,
	)

	req := connect.NewRequest(&authnv1.AuthenticateRequest{
		Token: testToken,
	})

	// act
	resp, err := handler.Authenticate(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "unauthenticated: missing or invalid authentication")
}

func Test_Handler_Authenticate_SubjectError(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	claims := authnrpc.NewMockClaims(t)
	claims.EXPECT().GetSubject().Return("", errors.New("error")).Once()

	parser := authnrpc.NewMockParser(t)
	parser.EXPECT().Parse(testToken).Return(claims, nil).Once()

	handler := authnrpc.New(
		authnrpc.NewMockIssuer(t),
		parser,
	)

	req := connect.NewRequest(&authnv1.AuthenticateRequest{
		Token: testToken,
	})

	// act
	resp, err := handler.Authenticate(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}
