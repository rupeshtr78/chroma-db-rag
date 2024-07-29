package langchain

import (
	"context"
	"errors"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
)

// Poc using langchain
func PdfToDocument(ctx context.Context, file string) ([]schema.Document, error) {
	// Open the pdf file
	f, err := os.Open("./test/Model Params.pdf")
	if err != nil {
		return nil, errors.New("could not open file")
	}
	defer f.Close()

	// Get the file info
	finfo, err := f.Stat()
	if err != nil {
		return nil, errors.New("could not stat file")
	}

	// Load the pdf file
	p := documentloaders.NewPDF(f, finfo.Size())
	docs, err := p.Load(ctx)
	if err != nil {
		return nil, errors.New("could not load pdf")
	}

	return docs, nil

}
