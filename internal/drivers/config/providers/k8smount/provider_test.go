package k8smount_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/drivers/config/providers/k8smount"
)

func createFile(t *testing.T, root, path, content string) {
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

	createFile(t, dir, "a", "a")
	createFile(t, dir, "b.c", "c")
	createFile(t, dir, "d.e.f", "f")

	provider := k8smount.Provider(dir, "." /*delim*/)

	// act
	values, err := provider.Read()

	// assert
	expected := map[string]any{
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

	assert.Equal(t, expected, values)
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
