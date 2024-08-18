package vectordb

import (
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
}

type ChromagoCollection struct {
	Collection *chromago.Collection
}

func NewChromagoCollection(collection *chromago.Collection) *ChromagoCollection {
	return &ChromagoCollection{Collection: collection}
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

type RecordSet interface {
	WithRecord(recordOpts ...types.Option) *types.RecordSet
}

type RecordSetWrapper struct {
	*types.RecordSet
}

func (rsw *RecordSetWrapper) WithRecord(recordOpts ...types.Option) *types.RecordSet {
	return rsw.RecordSet.WithRecord(recordOpts...)
}
