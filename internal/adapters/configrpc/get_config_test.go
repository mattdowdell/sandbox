package configrpc_test

import (
	"errors"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/gen/config/v1"
	"github.com/mattdowdell/sandbox/internal/adapters/configrpc"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockconfigrpc"
	"github.com/mattdowdell/sandbox/pkg/slogt"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Test_Handler_GetConfig_Success(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	loader := mockconfigrpc.NewLoader(t)
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

	req := connect.NewRequest(&configv1.GetConfigRequest{})

	// act
	resp, err := handler.GetConfig(ctx, req)

	// assert
	want := connect.NewResponse(&configv1.GetConfigResponse{
		Config: map[string]string{
			"foo": "foo",
		},
	})

	assert.Equal(t, want, resp)
	assert.NoError(t, err)
}

func Test_Handler_GetConfig_LoadError(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))

	loader := mockconfigrpc.NewLoader(t)
	loader.
		EXPECT().
		Load(mockConfig).
		Return(errors.New("example")).
		Once()

	handler := configrpc.New[TestConfig](loader)

	req := connect.NewRequest(&configv1.GetConfigRequest{})

	// act
	resp, err := handler.GetConfig(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}

func Test_Handler_GetConfig_EncodeError(t *testing.T) {
	// arrange
	ctx := slogx.IntoContext(t.Context(), slogt.New(t))
	val := new(bool)

	loader := mockconfigrpc.NewLoader(t)
	loader.EXPECT().Load(val).Return(nil).Once()

	handler := configrpc.New[bool](loader)

	req := connect.NewRequest(&configv1.GetConfigRequest{})

	// act
	resp, err := handler.GetConfig(ctx, req)

	// assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "internal: internal error")
}
