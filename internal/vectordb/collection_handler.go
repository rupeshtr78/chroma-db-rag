package vectordb

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

// EmbeddingFunc interface depedency injection for embedding documents
type EmbeddingFunc interface {
	EmbedDocuments(ctx context.Context, docs []string) ([]*types.Embedding, error)
}

// Collection interface allows dependency injection for the collection
type Collection interface {
	QueryWithOptions(ctx context.Context, options ...types.CollectionQueryOption) (*chromago.QueryResults, error)
	EmbeddingFunction() EmbeddingFunc
	AddRecords(ctx context.Context, recordSet *types.RecordSet) (*chromago.Collection, error)
	AddRecordSetToCollection(ctx context.Context, recordSet *types.RecordSet, docs []string, metadata constants.Metadata) (*chromago.Collection, error)
}

type ChromagoCollection struct {
	Collection *chromago.Collection
}

// NewChromagoCollection creates a new collection with the given name, embedding function and distance function.
func NewChromagoCollection(ctx context.Context, chromaClient *chromaclient.ChromaClient, embeddingFunction types.EmbeddingFunction) (*ChromagoCollection, error) {
	// CreateCollection creates a new collection with the given name, embedding function and distance function.
	collection, err := chromaclient.CreateCollection(ctx, chromaClient, embeddingFunction)
	return &ChromagoCollection{Collection: collection}, err
}

func (c *ChromagoCollection) EmbeddingFunction() EmbeddingFunc {
	return c.Collection.EmbeddingFunction
}

func (c *ChromagoCollection) AddRecords(ctx context.Context, recordSet *types.RecordSet) (*chromago.Collection, error) {
	cc, err := c.Collection.AddRecords(ctx, recordSet)
	return cc, err
}

func (c *ChromagoCollection) QueryWithOptions(ctx context.Context, options ...types.CollectionQueryOption) (*chromago.QueryResults, error) {
	qr, err := c.Collection.QueryWithOptions(ctx, options...)
	return qr, err
}
