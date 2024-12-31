package k8smount

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/v2"
)

// Non-allocating compile-time check for interface implementation.
var _ koanf.Provider = (*K8SMount)(nil)

// K8SMount implements a Kubernetes pod mount provider.
type K8SMount struct {
	mount string
	delim string
}

// Provider creates a new K8SMount provider capable of reading in mounted secrets and configmaps in
// a Kubernetes pod.
//
// The given path should be the mount point of the configmap or secret. The delimiter is used to
// create a heirarchy of keys based on the mounted filename. For example, a configmap mounted at
// "/my/config/" with a key of "log.level" set to "INFO" would result in {"log":{"level":"INFO"}}
// being read as configuration.
func Provider(mount, delim string) *K8SMount {
	return &K8SMount{
		mount: mount,
		delim: delim,
	}
}

// ReadBytes is not supported by the provider.
func (*K8SMount) ReadBytes() ([]byte, error) {
	return nil, errors.New("k8smount provider does not support this method")
}

// Read collects the contents of all files under the mount point and returns them as a map.
func (k *K8SMount) Read() (map[string]any, error) {
	values := map[string]any{}
	dirFs := os.DirFS(k.mount)

	if err := fs.WalkDir(dirFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip sub-directories as k8s resources can't use path separators in keys
		if d.IsDir() {
			if path == "." {
				return nil
			}

			return fs.SkipDir
		}

		file, err := dirFs.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		key := strings.TrimPrefix(path, k.mount)
		values[key] = string(content)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to read configuration from mount: %q: %w", k.mount, err)
	}

	return maps.Unflatten(values, k.delim), nil
}
