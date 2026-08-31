package provider

import (
	"fmt"
	"os"
	"sync"
)

type PromptLoader struct {
	cache map[string]string
	mutex sync.RWMutex
}

func NewPromptLoader() *PromptLoader {
	return &PromptLoader{
		cache: make(map[string]string),
	}
}

// Load reads a prompt from the given file path and caches it.
// Example paths:
//   - internal/resume/prompts/resume_parser_v1.txt
//   - internal/interview/engine/prompts/technical_question_v1.txt
//   - internal/evaluation/prompts/final_evaluation_v1.txt
func (p *PromptLoader) Load(promptPath string) (string, error) {

	p.mutex.RLock()
	if prompt, ok := p.cache[promptPath]; ok {
		p.mutex.RUnlock()
		return prompt, nil
	}
	p.mutex.RUnlock()

	content, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("load prompt %s: %w", promptPath, err)
	}

	prompt := string(content)

	p.mutex.Lock()
	p.cache[promptPath] = prompt
	p.mutex.Unlock()

	return prompt, nil
}