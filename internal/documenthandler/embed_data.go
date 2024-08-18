package documenthandler

import (
	"chroma-db/internal/constants"
	"chroma-db/internal/vectordb"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
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

// VectorEmbedData embeds the data in the collection
func VectorEmbedData(ctx context.Context, c vectordb.Collection, recordSet *types.RecordSet, options ...Option) (*chromago.Collection, error) {
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

	// collection, recordSet, err := vectordb.CreateCollectionAndRecordSet(ctx, client, constants.HuggingFace, opts.EmbeddingModel)
	// if err != nil {
	// 	log.Debug().Msgf("Error creating collection and recordset: %v\n", err)
	// 	return nil, err
	// }

	docLoader := NewDocumentLoader(opts.DocType)

	docs, metadata, err := docLoader.LoadDocument(ctx, opts.DocPath)
	// Load the documents
	// docs, metadata, err := LoadTextDocuments(ctx, docPath)
	if err != nil {
		log.Debug().Msgf("Error loading documents: %v\n", err)
		return nil, err
	}

	// Add the record set to the collection
	collection, err := c.AddRecordSetToCollection(ctx, recordSet, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding record set to collection: %v\n", err)
		return nil, err
	}

	// Count the number of documents in the collection
	countDocs, qrerr := collection.Count(ctx)
	if qrerr != nil {
		log.Debug().Msgf("Error counting documents: %s \n", qrerr)
	}

	if countDocs == 0 {
		log.Debug().Msgf("No documents found in the collection\n")
	}

	return collection, nil

}
