package k8smount

import (
	"errors"
	"fmt"
	"io"
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

var (
	ErrUnexpectedEvent = errors.New("unexpected event")
	ErrAlreadyWatched  = errors.New("mount is already being watched")
)

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
// The given path should be the mount point of the configmap or secret. The delimiter is used to
// create a hierarchy of keys based on the mounted filename. For example, a configmap mounted at
// "/my/config/" with a key of "log.level" set to "INFO" would result in {"log":{"level":"INFO"}}
// being read as configuration.
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

// Watch starts a watcher in a goroutine for the files under the mount point and calls the given
// function when changes occur.
//
// Only one watcher may be started at a time.
//
// If an error occurs, the function is called with the error before the watch is stopped. If the
// function is called with a nil error value, a change was detected successfully and watching will
// continue.
func (k *K8SMount) Watch(fn func(err error)) error {
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
func (k *K8SMount) watchDir(fn func(err error)) {
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

			fmt.Println(event.String())

			// mounts are only meant to be updated for the lifetime of the pod
			if !event.Has(fsnotify.Write | fsnotify.Chmod) {
				fn(fmt.Errorf("%w: %s", ErrUnexpectedEvent, event.String()))
				return
			}

			// ignore chmod
			if !event.Has(fsnotify.Chmod) {
				fn(nil)
			}

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
