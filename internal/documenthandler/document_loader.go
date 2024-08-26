package documenthandler

import (
	"chroma-db/internal/constants"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

type DocumentLoader interface {
	LoadDocument(context.Context, string) ([]string, constants.Metadata, error)
}

// Implementing the file reader interface for real file operations
type FileReader struct{}

func (r *FileReader) ReadFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

// NewDocumentLoader returns a new DocumentLoader based on the docType
func NewDocumentLoader(docType constants.DocType) DocumentLoader {
	switch docType {
	case constants.PDF:
		return &PdfLoader{}
	case constants.TXT:
		return &TextLoader{}
	default:
		return nil
	}
}

type TextLoader struct {
	fileReader FileReader
}

// LoadDocument loads the text data from the given file path
func (t *TextLoader) LoadDocument(ctx context.Context, filePath string) ([]string, constants.Metadata, error) {
	if filePath == "" {
		return nil, nil, fmt.Errorf("TextLoaderV2: File path is empty")
	}
	f, err := t.fileReader.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	loader := documentloaders.NewText(f)
	// docs, err := loader.Load(ctx)

	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(1024),
		textsplitter.WithChunkOverlap(512),
	)
	docs, err := loader.LoadAndSplit(ctx, splitter)

	if err != nil || len(docs) == 0 {
		return nil, nil, err
	}

	// returning a slice of strings and a Metadata map
	strSlice := make([]string, 0)
	metadata := make(constants.Metadata)
	for _, doc := range docs {
		strSlice = append(strSlice, doc.PageContent)
		if doc.Metadata == nil {
			continue
		}
		for k, v := range doc.Metadata {
			log.Debug().Msgf("TextLoaderV2: Metadata: %s: %s", k, v)
			metadata[k] = v
		}
	}

	log.Debug().Msgf("TextLoader: Loading loaded text data from %s", filePath)
	log.Debug().Msgf("TextLoader: Metadata: %v", metadata)

	return strSlice, metadata, nil
}

type PdfLoader struct{}

// LoadDocument loads the text data from the given pdf file path / TODO - add metadata fix issues with split
func (p *PdfLoader) LoadDocument(ctx context.Context, filePath string) ([]string, constants.Metadata, error) {
	if filePath == "" {
		return nil, nil, fmt.Errorf("TextLoaderV2: File path is empty")
	}
	pdfStrings := []string{}
	metadata := constants.Metadata{}

	file, pdfReader, err := pdf.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	fileName := strings.Split(file.Name(), "/")[1]

	numPages := pdfReader.NumPage()
	log.Debug().Msgf("Number of pdf pages: %d", numPages)

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page := pdfReader.Page(pageNum)

		text, err := page.GetPlainText(nil)
		if err != nil {
			return nil, nil, err
		}
		pdfStrings = append(pdfStrings, text)
		metadata[fmt.Sprintf("%d", pageNum)] = fileName

	}

	log.Debug().Msgf("PDF List Length: %v", len(pdfStrings))
	log.Debug().Msgf("Metadata Length: %v", len(metadata))

	return pdfStrings, metadata, nil
}
