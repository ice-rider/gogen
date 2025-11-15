package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Writer struct {
	projectRoot string
	written     []string
	mu          sync.Mutex
}

func NewWriter(projectRoot string) *Writer {
	return &Writer{
		projectRoot: projectRoot,
		written:     make([]string, 0),
	}
}

func (w *Writer) Write(relativePath, content string, overwrite bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	fullPath := filepath.Join(w.projectRoot, relativePath)

	if _, err := os.Stat(fullPath); err == nil && !overwrite {
		return fmt.Errorf("file already exists: %s", relativePath)
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	w.written = append(w.written, fullPath)

	return nil
}

func (w *Writer) WriteIfNotExists(relativePath, content string) error {
	return w.Write(relativePath, content, false)
}

func (w *Writer) Rollback() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var errors []error

	for i := len(w.written) - 1; i >= 0; i-- {
		path := w.written[i]

		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("failed to remove %s: %w", path, err))
		}
	}

	w.written = make([]string, 0)

	if len(errors) > 0 {
		return fmt.Errorf("rollback completed with errors: %v", errors)
	}

	return nil
}

func (w *Writer) GetWrittenFiles() []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	result := make([]string, len(w.written))
	copy(result, w.written)
	return result
}

func (w *Writer) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.written = make([]string, 0)
}
