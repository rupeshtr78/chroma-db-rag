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
	// AddRecordSetToCollection(ctx context.Context, recordSet *types.RecordSet, docs []string, metadata map[string]interface{}) (*chromago.Collection, error)
}

type ChromagoCollection struct {
	*chromago.Collection
}

func (ccc *ChromagoCollection) EmbeddingFunction() EmbeddingFunc {
	return ccc.Collection.EmbeddingFunction
}
