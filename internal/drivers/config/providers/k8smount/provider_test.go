package k8smount_test

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/drivers/config/providers/k8smount"
)

func Test_New(t *testing.T) {
	// arrange
	dir := t.TempDir()

	// act
	provider := k8smount.Provider(dir, "." /*delim*/)

	// assert
	assert.NotNil(t, provider)
}

func Test_K8SMount_ReadBytes(t *testing.T) {
	// arrange
	dir := t.TempDir()
	provider := k8smount.Provider(dir, "." /*delim*/)

	// act
	content, err := provider.ReadBytes()

	// assert
	assert.Empty(t, content)
	assert.EqualError(t, err, "k8smount provider does not support this method")
}

func Test_K8SMount_Read_Empty(t *testing.T) {
	// arrange
	dir := t.TempDir()
	provider := k8smount.Provider(dir, "." /*delim*/)

	// act
	values, err := provider.Read()

	// assert
	assert.Empty(t, values)
	assert.NoError(t, err)
}

func Test_K8SMount_Read_WithFiles(t *testing.T) {
	// arrange
	dir := t.TempDir()

	require.NoError(t, writeFile(filepath.Join(dir, "a"), "a"))
	require.NoError(t, writeFile(filepath.Join(dir, "b.c"), "c"))
	require.NoError(t, writeFile(filepath.Join(dir, "d.e.f"), "f"))
	require.NoError(t, writeFile(filepath.Join(dir, "dir", "file"), "a")) // should be ignored

	provider := k8smount.Provider(dir, "." /*delim*/)

	// act
	got, err := provider.Read()

	// assert
	want := map[string]any{
		"a": "a",
		"b": map[string]any{
			"c": "c",
		},
		"d": map[string]any{
			"e": map[string]any{
				"f": "f",
			},
		},
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_K8SMount_Read_WithVolumeMount(t *testing.T) {
	// arrange
	dir := t.TempDir()

	require.NoError(t, writeVolumeMount(dir, map[string]string{
		"a.foo": "foo-value",
		"b.bar": "bar-value",
		"b.baz": "baz-value",
	}))

	provider := k8smount.Provider(dir, "." /*delim*/)

	// act
	got, err := provider.Read()

	// assert
	want := map[string]any{
		"a": map[string]any{
			"foo": "foo-value",
		},
		"b": map[string]any{
			"bar": "bar-value",
			"baz": "baz-value",
		},
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_K8SMount_Read_MissingDir(t *testing.T) {
	// arrange
	provider := k8smount.Provider("/does/not/exist" /*dir*/, "." /*delim*/)

	// act
	values, err := provider.Read()

	// assert
	assert.Empty(t, values)
	assert.EqualError(
		t,
		err,
		`failed to open mount: open /does/not/exist: no such file or directory`,
	)
}

func Test_K8SMount_Watch_Success(t *testing.T) {
	// arrange
	dir := t.TempDir()

	require.NoError(t, writeFile(filepath.Join(dir, "a"), "a"))

	provider := k8smount.Provider(dir, "." /*delim*/)

	_, err := provider.Read()
	require.NoError(t, err)

	var watched atomic.Bool

	// act
	require.NoError(t, provider.Watch(func(err error) {
		assert.NoError(t, err)
		watched.Store(true)
	}))

	for !watched.Load() {
		require.NoError(t, writeFile(filepath.Join(dir, "a"), "b"))
	}

	// assert
	require.NoError(t, provider.Unwatch())

	got, err := provider.Read()

	want := map[string]any{
		"a": "b",
	}

	assert.Equal(t, want, got)
	assert.NoError(t, err)
}

func Test_K8SMount_Watch_AlreadyWatching(t *testing.T) {
	// arrange
	dir := t.TempDir()
	provider := k8smount.Provider(dir, "." /*delim*/)

	require.NoError(t, provider.Watch(func(err error) {
		assert.NoError(t, err)
	}))
	defer func() {
		assert.NoError(t, provider.Unwatch())
	}()

	// act
	err := provider.Watch(func(err error) {
		assert.NoError(t, err)
	})

	// assert
	assert.ErrorIs(t, err, k8smount.ErrAlreadyWatched)
}

func Test_K8SMount_Watch_UnexpectedEvent(t *testing.T) {
	// arrange
	dir := t.TempDir()

	require.NoError(t, writeFile(filepath.Join(dir, "a"), "a"))

	provider := k8smount.Provider(dir, "." /*delim*/)

	_, err := provider.Read()
	require.NoError(t, err)

	watched := make(chan struct{})

	// act
	require.NoError(t, provider.Watch(func(err error) {
		assert.ErrorIs(t, err, k8smount.ErrUnexpectedEvent)
		close(watched)
	}))

	os.Remove(filepath.Join(dir, "a"))

	// assert
	<-watched
	require.NoError(t, provider.Unwatch())
}

func Test_K8SMount_Unwatch(t *testing.T) {
	// arrange
	dir := t.TempDir()
	provider := k8smount.Provider(dir, "." /*delim*/)

	// act
	err := provider.Unwatch()

	// assert
	assert.NoError(t, err)
}
