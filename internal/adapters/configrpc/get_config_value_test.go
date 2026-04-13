package configrpc_test

import (
	"errors"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mattdowdell/sandbox/gen/config/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/configrpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

var mockConfig = mock.AnythingOfType("*configrpc_test.TestConfig")

func Test_Handler_GetConfigValue_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	loader := configrpc.NewMockLoader(t)
	loader.
		EXPECT().
		Load(mockConfig).
		RunAndReturn(func(target any) error {
			tgt, ok := target.(*TestConfig)
			if !ok {
				return fmt.Errorf("unexpected target type: %T", target)
			}

			*tgt = TestConfig{Foo: "foo"}
			return nil
		}).
		Once()

	handler := configrpc.New[TestConfig](loader)

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

	loader := configrpc.NewMockLoader(t)
	loader.
		EXPECT().
		Load(mockConfig).
		Return(errors.New("example")).
		Once()

	handler := configrpc.New[TestConfig](loader)

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
	val := new(bool)

	loader := configrpc.NewMockLoader(t)
	loader.EXPECT().Load(val).Return(nil).Once()

	handler := configrpc.New[bool](loader)

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

	loader := configrpc.NewMockLoader(t)
	loader.
		EXPECT().
		Load(mockConfig).
		RunAndReturn(func(target any) error {
			tgt, ok := target.(*TestConfig)
			if !ok {
				return fmt.Errorf("unexpected target type: %T", target)
			}

			*tgt = TestConfig{Foo: "foo"}
			return nil
		}).
		Once()

	handler := configrpc.New[TestConfig](loader)

	req := connect.NewRequest(&configv1.GetConfigValueRequest{
		Key: "bar",
	})

	// act
	resp, err := handler.GetConfigValue(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "not_found: value does not exist")
}
