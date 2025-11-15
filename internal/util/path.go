package util

import (
	"path/filepath"
	"strings"
)

func NormalizePath(path string) string {
	return filepath.ToSlash(path)
}

func JoinModulePath(modulePath, relativePath string) string {
	normalized := NormalizePath(relativePath)
	normalized = strings.TrimPrefix(normalized, "/")
	normalized = strings.TrimSuffix(normalized, ".go")

	return modulePath + "/" + normalized
}

func GetPackageName(path string) string {
	normalized := NormalizePath(path)
	parts := strings.Split(normalized, "/")

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}
