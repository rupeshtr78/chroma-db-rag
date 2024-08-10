package ollamarag

import (
	"chroma-db/app/ollamarag"
	"chroma-db/app/vectordb"
	"chroma-db/internal/constants"
	"chroma-db/internal/documenthandler"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/rs/zerolog/log"
)

// Option is a type for argument options
type Option func(*ollamaRagOptions)

// ollamaRagOptions is the options for RunOllamaRag
type ollamaRagOptions struct {
	ChromaURL      string
	TenantName     string
	DatabaseName   string
	DocPath        string
	DocType        constants.DocType
	EmbeddingModel string
}

// WithChromaURL sets the Chroma URL
func WithChromaURL(url string) Option {
	return func(o *ollamaRagOptions) {
		o.ChromaURL = url
	}
}

// WithTenantName sets the tenant name
func WithTenantName(name string) Option {
	return func(o *ollamaRagOptions) {
		o.TenantName = name
	}
}

// WithDatabaseName sets the database name
func WithDatabaseName(name string) Option {
	return func(o *ollamaRagOptions) {
		o.DatabaseName = name
	}
}

// WithDocPath sets the document path
func WithDocPath(path string) Option {
	return func(o *ollamaRagOptions) {
		o.DocPath = path
	}
}

// WithDocType sets the document type
func WithDocType(docType constants.DocType) Option {
	return func(o *ollamaRagOptions) {
		o.DocType = docType
	}
}

// WithEmbeddingModel sets the embedding model
func WithEmbeddingModel(model string) Option {
	return func(o *ollamaRagOptions) {
		o.EmbeddingModel = model
	}
}

// RunOllamaRag runs the Ollama RAG process with the given options
// ChromaURL:      constants.ChromaUrl,
// TenantName:     constants.TenantName,
// DatabaseName:   constants.Database,
// EmbeddingModel: constants.OllamaEmbdedModel,
func RunOllamaRagV2(ctx context.Context, options ...Option) (*chromago.Collection, error) {
	// Default options
	opts := &ollamaRagOptions{
		ChromaURL:      constants.ChromaUrl,
		TenantName:     constants.TenantName,
		DatabaseName:   constants.Database,
		DocPath:        "default-path",
		DocType:        constants.TXT,
		EmbeddingModel: constants.OllamaEmbdedModel,
	}

	// Apply the options
	for _, option := range options {
		option(opts)
	}

	// Initialize the Chroma client
	collection, recordSet, err := vectordb.InitializeChroma(ctx, opts.ChromaURL, opts.TenantName, opts.DatabaseName, opts.EmbeddingModel)
	if err != nil {
		log.Debug().Msgf("Error initializing Chroma: %v\n", err)
		return nil, err
	}

	docLoader := documenthandler.NewDocumentLoader(opts.DocType)

	docs, metadata, err := docLoader.LoadDocument(ctx, opts.DocPath)
	// Load the documents
	// docs, metadata, err := LoadTextDocuments(ctx, docPath)
	if err != nil {
		log.Debug().Msgf("Error loading documents: %v\n", err)
		return nil, err
	}

	// Add the record set to the collection
	collection, err = ollamarag.AddRecordSetToCollection(ctx, collection, recordSet, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding record set to collection: %v\n", err)
		return nil, err
	}

	// Count the number of documents in the collection
	countDocs, qrerr := collection.Count(ctx)
	if qrerr != nil {
		log.Debug().Msgf("Error counting documents: %s \n", qrerr)
	}

	log.Debug().Msgf("Number of documents in collection: %d\n", countDocs)

	return collection, nil

}
