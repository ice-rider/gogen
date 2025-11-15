package template

import (
	"bytes"
	"fmt"
)

type Renderer struct {
	loader *Loader
}

func NewRenderer(loader *Loader) *Renderer {
	return &Renderer{
		loader: loader,
	}
}

func (r *Renderer) Render(templateName string, data interface{}) (string, error) {

	tmpl, err := r.loader.Load(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func (r *Renderer) RenderToFile(templateName string, data interface{}, outputPath string) error {
	content, err := r.Render(templateName, data)
	if err != nil {
		return err
	}

	_ = content
	_ = outputPath

	return nil
}
