package vectordb

import (
	"chroma-db/internal/constants"
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/rs/zerolog/log"
)

// AddRecordSetToCollection adds the record set to the collection.
func AddRecordSetToCollection(ctx context.Context,
	collection *chromago.Collection,
	recordSet *types.RecordSet,
	docs []string,
	metadata constants.Metadata) (*chromago.Collection, error) {
	// Add the documents to the record set
	recordSet, err := AddTextToRecordSet(ctx, recordSet, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding to record set: %v\n", err)
		return nil, err
	}

	// Build and validate the record set
	_, err = recordSet.BuildAndValidate(ctx)
	if err != nil {
		log.Debug().Msgf("Error building and validating records: %v\n", err)
		return nil, err
	}

	// Add the records to the collection
	collection, err = collection.AddRecords(ctx, recordSet)
	if err != nil {
		log.Err(err).Msgf("Error adding records to collection: %s\n", collection.Name)
		return nil, err
	}

	log.Debug().Msgf("Added %d records to collection: %s\n", len(docs), collection.Name)

	return collection, nil
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
