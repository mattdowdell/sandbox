// Package k8smount contains a Koanf provider for loading configuration from Kubernetes volume
// mounts, i.e. Secrets or ConfigMaps mounted into a Pod.
//
// This is most appropriate for key-value data. Using a ConfigMap as an example:
//
//	apiVersion: v1
//	kind: ConfigMap
//	metadata:
//	  name: example
//	data:
//	  foo: "true"
//	  bar: "1"
//	  baz: "value"
//
// If ConfigMap fields contains structured data, such as JSON or YAML, the file provider with the
// appropriate parser should be used instead.
package k8smount

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/v2"
)

var ErrAlreadyWatched = errors.New("mount is already being watched")

// Non-allocating compile-time check for interface implementation.
var _ koanf.Provider = (*K8SMount)(nil)

// K8SMount implements a Kubernetes pod mount provider.
type K8SMount struct {
	mount    string
	delim    string
	watching atomic.Bool
	watcher  *fsnotify.Watcher
}

// Provider creates a new K8SMount provider capable of reading in mounted secrets and configmaps in
// a Kubernetes pod.
//
// The given mount should be the mount point of the configmap or secret. The delimiter is used to
// create a hierarchy of keys based on the mounted filename. For example, a configmap mounted at
// "/mnt/config/" with a key of "log.level" set to "INFO" and a delimiter of "." would result in
// {"log":{"level":"INFO"}} being read as configuration.
//
// Keys mounted in directories are always split. For example, if the above key was mounted at
// "log/level" instead, it will always produce {"log":{"level":"INFO"}} as the result.
func Provider(mount, delim string) *K8SMount {
	return &K8SMount{
		mount: filepath.Clean(mount),
		delim: delim,
	}
}

// ReadBytes is not supported by the provider.
func (*K8SMount) ReadBytes() ([]byte, error) {
	return nil, errors.New("k8smount provider does not support this method")
}

// Read collects the contents of all files under the mount point and returns them as a map.
func (k *K8SMount) Read() (map[string]any, error) {
	root, err := os.OpenRoot(k.mount)
	if err != nil {
		return nil, fmt.Errorf("failed to open mount: %w", err)
	}

	data := map[string]any{}
	mountFS := root.FS()

	if err := fs.WalkDir(mountFS, ".", func(path string, d fs.DirEntry, err error) error {
		key, content, err := k.walkDir(mountFS, path, d, err)
		if err != nil {
			return err
		}

		if key != "" {
			data[key] = content
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to read configuration from mount: %q: %w", k.mount, err)
	}

	return maps.Unflatten(maps.Unflatten(data, k.delim), "/"), nil
}

//nolint:gocritic // named return vars don't make code much more readable
func (k *K8SMount) walkDir(mountFS fs.FS, path string, d fs.DirEntry, err error) (string, string, error) {
	if err != nil {
		return "", "", err
	}

	if path == "." && d.IsDir() {
		return "", "", nil
	}

	resolved := path

	for d.Type()&os.ModeSymlink != 0 {
		p, err := fs.ReadLink(mountFS, resolved)
		if err != nil {
			return "", "", err
		}

		info, err := fs.Lstat(mountFS, p)
		if err != nil {
			return "", "", err
		}

		d = fs.FileInfoToDirEntry(info)
		resolved = p
	}

	if d.IsDir() {
		// don't skip ..data as it causes the symlinks we want to look at to be skipped
		// instead, just skip the timestamp based directories
		if strings.HasPrefix(path, "..") && path != "..data" {
			return "", "", fs.SkipDir
		}

		return "", "", nil
	}

	content, err := fs.ReadFile(mountFS, path)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}

	key := strings.TrimPrefix(path, k.mount)
	return key, string(content), nil
}

// Watch starts a watcher in a goroutine for the files under the mount point and calls the given
// function when changes occur.
//
// Only one watcher may be started at a time.
//
// If an error occurs, the function is called with the error before the watch is stopped. If the
// function is called with a nil error value, a change was detected successfully and watching will
// continue.
func (k *K8SMount) Watch(fn func(error)) error {
	if k.watching.Swap(true) {
		return ErrAlreadyWatched
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	k.watcher = watcher
	go k.watchDir(fn)

	return watcher.Add(k.mount)
}

//nolint:gocognit // short enough that the complaxity is acceptable
func (k *K8SMount) watchDir(fn func(error)) {
	defer k.watching.Store(false)

	var (
		lastEvent     string
		lastEventTime time.Time
	)

	for {
		select {
		case event, ok := <-k.watcher.Events:
			if !ok {
				return
			}

			// Use a simple timer to buffer events as certain events fire
			// multiple times on some platforms.
			if event.String() == lastEvent && time.Since(lastEventTime) < time.Millisecond*5 {
				continue
			}

			lastEvent = event.String()
			lastEventTime = time.Now()

			// ignore chmod
			if event.Has(fsnotify.Chmod) {
				return
			}

			// TODO: report which key was updated
			fn(nil)

		case err, ok := <-k.watcher.Errors:
			if !ok {
				return
			}

			fn(err)
			return
		}
	}
}

// Unwatch stops a previously started Watch.
func (k *K8SMount) Unwatch() error {
	if k.watcher != nil {
		return k.watcher.Close()
	}

	return nil
}
