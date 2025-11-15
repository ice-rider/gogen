package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (f *Finder) GetModulePath() (string, error) {
	root, err := f.FindRoot()
	if err != nil {
		return "", err
	}

	goModPath := filepath.Join(root, "go.mod")

	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			return modulePath, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	return "", fmt.Errorf("module path not found in go.mod")
}

func (f *Finder) GetModuleInfo() (root, modulePath string, err error) {
	root, err = f.FindRoot()
	if err != nil {
		return "", "", err
	}

	modulePath, err = f.GetModulePath()
	if err != nil {
		return "", "", err
	}

	return root, modulePath, nil
}
