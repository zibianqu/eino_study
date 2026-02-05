package utils

import (
	"os"
	"path/filepath"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetFileExtension returns the file extension
func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

// GetFileName returns the file name without extension
func GetFileName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return base[:len(base)-len(ext)]
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}