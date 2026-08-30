package resume

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Extractor interface {
	ExtractText(file multipart.File) (string, error)
}

type extractor struct{}

func NewExtractor() Extractor {
	return &extractor{}
}

func (e *extractor) ExtractText(file multipart.File) (string, error) {

	// Read uploaded PDF into memory.
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read uploaded pdf: %w", err)
	}

	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("open pdf reader: %w", err)
	}

	var builder strings.Builder

	totalPages := reader.NumPage()

	for pageNumber := 1; pageNumber <= totalPages; pageNumber++ {

		page := reader.Page(pageNumber)

		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("extract page %d text: %w", pageNumber, err)
		}

		builder.WriteString(text)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}