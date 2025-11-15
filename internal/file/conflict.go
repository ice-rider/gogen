package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ConflictResolver struct {
	interactive    bool
	forceOverwrite bool
}

func NewConflictResolver(interactive, forceOverwrite bool) *ConflictResolver {
	return &ConflictResolver{
		interactive:    interactive,
		forceOverwrite: forceOverwrite,
	}
}

func (cr *ConflictResolver) ResolveConflict(path string) (bool, error) {

	if cr.forceOverwrite {
		return true, nil
	}

	if !cr.interactive {
		return false, fmt.Errorf("file exists: %s (use --force to overwrite)", path)
	}

	fmt.Printf("File %s already exists. Overwrite? (y/N): ", path)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y" || response == "yes", nil
}

func (cr *ConflictResolver) CheckConflicts(paths []string) ([]string, error) {
	var conflicts []string

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			conflicts = append(conflicts, path)
		}
	}

	return conflicts, nil
}
