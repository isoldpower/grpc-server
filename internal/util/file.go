package util

import (
	"os"
	"path/filepath"
	"runtime"
)

// ResolvePath resolves the specified path based on the base or current working directory
func ResolvePath(path string, base string) string {
	if filepath.IsAbs(path) {
		return path
	}

	if path == "" {
		path = "."
	}

	if cwd, err := os.Getwd(); err == nil && base == "" {
		base = cwd
	}

	return filepath.Join(base, path)
}

// RelativePath returns the path relative to base or current working directory
func RelativePath(path string, base string) string {
	if path == "" {
		path = "."
	}

	if cwd, err := os.Getwd(); base == "" && err == nil {
		base = cwd
	}

	if rel, err := filepath.Rel(base, path); err == nil {
		return rel
	}

	return path
}

// RootPath returns project's root directory absolute path
func RootPath() string {
	_, file, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(file), "../..")
}
