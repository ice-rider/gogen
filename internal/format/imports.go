package format

import (
	"bytes"
	"fmt"

	"golang.org/x/tools/imports"
)

type ImportsManager struct{}

func NewImportsManager() *ImportsManager {
	return &ImportsManager{}
}

func (im *ImportsManager) OrganizeImports(source string) (string, error) {
	organized, err := imports.Process("", []byte(source), nil)
	if err != nil {
		return "", fmt.Errorf("failed to organize imports: %w", err)
	}

	return string(organized), nil
}

func (im *ImportsManager) AddImports(source string, imports []string) string {

	lines := bytes.Split([]byte(source), []byte("\n"))

	var result bytes.Buffer
	importInserted := false

	for i, line := range lines {

		if !importInserted && bytes.HasPrefix(line, []byte("package ")) {
			result.Write(line)
			result.WriteByte('\n')

			if i+1 < len(lines) && len(bytes.TrimSpace(lines[i+1])) == 0 {
				result.WriteByte('\n')
				continue
			}

			if len(imports) > 0 {
				result.WriteString("\nimport (\n")
				for _, imp := range imports {
					result.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
				}
				result.WriteString(")\n")
			}

			importInserted = true
			continue
		}

		result.Write(line)
		result.WriteByte('\n')
	}

	return result.String()
}
