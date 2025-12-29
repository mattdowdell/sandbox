package configrpc_test

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/gen/config/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/configrpc"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockconfigrpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_Handler_GetConfigValue_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	loader := mockconfigrpc.NewLoader[TestConfig](t)
	loader.EXPECT().Load().Return(&TestConfig{Foo: "foo"}, nil).Once()

	handler := configrpc.New(loader)

	req := connect.NewRequest(&configv1.GetConfigValueRequest{
		Key: "foo",
	})

	// act
	resp, err := handler.GetConfigValue(ctx, req)

	// assert
	want := connect.NewResponse(&configv1.GetConfigValueResponse{
		Value: "foo",
	})

	assert.Equal(t, want, resp)
	assert.NoError(t, err)
}

func Test_Handler_GetConfigValue_LoadError(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	loader := mockconfigrpc.NewLoader[TestConfig](t)
	loader.EXPECT().Load().Return(nil, errors.New("example")).Once()

	handler := configrpc.New(loader)

	req := connect.NewRequest(&configv1.GetConfigValueRequest{})

	// act
	resp, err := handler.GetConfigValue(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}

func Test_Handler_GetConfigValue_EncodeError(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))
	val := true

	loader := mockconfigrpc.NewLoader[bool](t)
	loader.EXPECT().Load().Return(&val, nil).Once()

	handler := configrpc.New(loader)

	req := connect.NewRequest(&configv1.GetConfigValueRequest{
		Key: "foo",
	})

	// act
	resp, err := handler.GetConfigValue(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}

func Test_Handler_GetConfigValue_NotFound(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	loader := mockconfigrpc.NewLoader[TestConfig](t)
	loader.EXPECT().Load().Return(&TestConfig{Foo: "foo"}, nil).Once()

	handler := configrpc.New(loader)

	req := connect.NewRequest(&configv1.GetConfigValueRequest{
		Key: "bar",
	})

	// act
	resp, err := handler.GetConfigValue(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "not_found: value does not exist")
}
