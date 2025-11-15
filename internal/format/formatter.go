package format

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format(source string) (string, error) {

	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("syntax error: %w", err)
	}

	formatted, err := format.Source([]byte(source))
	if err != nil {
		return "", fmt.Errorf("format error: %w", err)
	}

	return string(formatted), nil
}

func (f *Formatter) FormatFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	formatted, err := f.Format(string(content))
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(formatted), 0644)
}

func (f *Formatter) Validate(source string) error {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	return err
}
