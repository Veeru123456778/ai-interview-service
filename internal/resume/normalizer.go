package resume

import (
	"regexp"
	"strings"
)

type Normalizer interface {
	Normalize(text string) string
}

type normalizer struct{}

func NewNormalizer() Normalizer {
	return &normalizer{}
}

func (n *normalizer) Normalize(text string) string {

	// Normalize line endings.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Replace tabs with spaces.
	text = strings.ReplaceAll(text, "\t", " ")

	// Remove multiple spaces.
	multipleSpaces := regexp.MustCompile(` +`)
	text = multipleSpaces.ReplaceAllString(text, " ")

	// Remove repeated blank lines.
	multipleNewLines := regexp.MustCompile(`\n{3,}`)
	text = multipleNewLines.ReplaceAllString(text, "\n\n")

	// Trim each line.
	lines := strings.Split(text, "\n")

	cleanLines := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines created by trimming.
		if line == "" {
			cleanLines = append(cleanLines, "")
			continue
		}

		cleanLines = append(cleanLines, line)
	}

	text = strings.Join(cleanLines, "\n")

	// Final trim.
	return strings.TrimSpace(text)
}