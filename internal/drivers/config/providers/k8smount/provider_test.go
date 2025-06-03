package k8smount_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/drivers/config/providers/k8smount"
)

func writeFile(t *testing.T, root, path, content string) {
	t.Helper()

	absPath := filepath.Join(root, path)
	absDirPath := filepath.Dir(absPath)

	if err := os.MkdirAll(absDirPath, 0o755); err != nil {
		t.Fatalf("failed to create directory: %s", err)
	}

	if err := os.WriteFile(absPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to create file: %s", err)
	}
}

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

	writeFile(t, dir, "a", "a")
	writeFile(t, dir, "b.c", "c")
	writeFile(t, dir, "d.e.f", "f")
	writeFile(t, dir, "dir/file", "a") // should be ignored

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
		`failed to read configuration from mount: "/does/not/exist": stat .: no such file or directory`,
	)
}

func Test_K8SMount_Watch_Success(t *testing.T) {
	// arrange
	dir := t.TempDir()

	writeFile(t, dir, "a", "a")

	provider := k8smount.Provider(dir, "." /*delim*/)

	_, err := provider.Read()
	require.NoError(t, err)

	watched := make(chan struct{})

	// act
	require.NoError(t, provider.Watch(func(err error) {
		assert.NoError(t, err)
		close(watched)
	}))

	writeFile(t, dir, "a", "b")

	// assert
	<-watched
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

	writeFile(t, dir, "a", "a")

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
