package chromaclient

import (
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

func CreateRecordSet(embeddingFunction types.EmbeddingFunction) (*types.RecordSet, error) {
	// Create a new record set with to hold the records to insert
	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(embeddingFunction),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Err(err).Msg("Error creating record set")
		return nil, err
	}

	return rs, nil
}

// AddPdfToRecordSet adds pdf documents to the record set
func AddPdfToRecordSet(ctx context.Context,
	collection *chromago.Collection,
	rs *types.RecordSet,
	documents []string,
	metadata map[string]any) (*types.RecordSet, error) {

	// Iterate over documents and metadata list and add records to the record set
	for i, doc := range documents {
		pageNum := i + 1
		key := fmt.Sprintf("%d", pageNum)
		metadataValue := metadata[key].(string)
		rs.WithRecord(
			types.WithDocument(doc),
			types.WithMetadata(key, metadataValue),
		)
	}

	return rs, nil
}

// internal/chromaclient/chroma_recordset.go
func AddTextToRecordSet(ctx context.Context,
	collection *chromago.Collection,
	rs *types.RecordSet,
	documents []string,
	metadata map[string]any) (*types.RecordSet, error) {

	// Iterate over documents and metadata list and add records to the record set
	for _, doc := range documents {
		rs.WithRecord(
			types.WithDocument(doc),
			// types.WithMetadatas(metadata),
		)
	}

	return rs, nil
}
