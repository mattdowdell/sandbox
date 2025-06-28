package k8smount_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	configmapTimeFmt = "2006_01_02_15_04_05.0000000000"
	dataDir          = "..data"
)

// writeVolumeMount creates a file structure that matches how a ConfigMap or Secret will be mounted
// in a Kubernetes Pod.
//
// First, files are created for each data field. These exist within a timestamp-based directory,
// likely when the ConfigMap or Secret was last modified.
//
//	..2025_06_28_09_28_32.3151791122/
//	├── database.hostname
//	├── database.name
//	└── database.port
//
// A symlink is then created for "..data" to the directory containing the data:
//
//	..data -> ..2025_06_28_09_28_32.3151791122
//
// Finally, symlinks are created for the data files, via the "..data" symlink:
//
//	database.hostname -> ..data/database.hostname
//	database.name -> ..data/database.name
//	database.port -> ..data/database.port
func writeVolumeMount(mount string, data map[string]string) error {
	dir := time.Now().UTC().Format(configmapTimeFmt)

	dirPath := filepath.Join(mount, dir)

	for key, value := range data {
		if err := writeFile(filepath.Join(dirPath, key), value); err != nil {
			return err
		}
	}

	if err := os.Chdir(mount); err != nil {
		return fmt.Errorf("failed to change dir to %q: %w", mount, err)
	}

	if err := os.Symlink(dir, "..data"); err != nil {
		return fmt.Errorf("failed to create %s symlink to %q: %w", "..data", dir, err)
	}

	for key := range data {
		target := filepath.Join(dataDir, key)

		if err := os.Symlink(target, key); err != nil {
			return fmt.Errorf("failed to create %q symlink to %q: %w", key, target, err)
		}
	}

	return nil
}

func writeFile(path, content string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to create file %q: %w", path, err)
	}

	fmt.Println("created:", path)

	return nil
}
