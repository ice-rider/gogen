package project

import (
	"fmt"
	"os"
	"path/filepath"
)

type Finder struct {
	startDir string
}

func NewFinder(startDir string) *Finder {
	if startDir == "" {
		startDir, _ = os.Getwd()
	}
	return &Finder{
		startDir: startDir,
	}
}

func (f *Finder) FindRoot() (string, error) {
	dir := f.startDir

	for {
		goModPath := filepath.Join(dir, "go.mod")

		if _, err := os.Stat(goModPath); err == nil {

			return dir, nil
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "", fmt.Errorf("go.mod not found: not a Go module")
		}

		dir = parent
	}
}

func (f *Finder) MustFindRoot() string {
	root, err := f.FindRoot()
	if err != nil {
		panic(err)
	}
	return root
}
