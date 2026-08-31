package prompts

import (
	"embed"
	"fmt"
)

//go:embed *.txt
var promptFS embed.FS

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(promptName string) (string, error) {

	fileName, ok := File(promptName)
	if !ok {
		return "", fmt.Errorf("prompt not registered: %s", promptName)
	}

	content, err := promptFS.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("load prompt %s: %w", promptName, err)
	}

	return string(content), nil
}
