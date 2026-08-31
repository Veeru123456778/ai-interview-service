package prompts

import (
	"bytes"
	"fmt"
	"text/template"
)

type Builder struct {
	loader *Loader
}

func NewBuilder(loader *Loader) *Builder {
	return &Builder{
		loader: loader,
	}
}

func (b *Builder) Build(
	promptName string,
	input map[string]any,
) (string, error) {

	templateText, err := b.loader.Load(promptName)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(promptName).Option("missingkey=error").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("parse prompt template: %w", err)
	}

	var output bytes.Buffer

	if err := tmpl.Execute(&output, input); err != nil {
		return "", fmt.Errorf("execute prompt template: %w", err)
	}

	return output.String(), nil
}