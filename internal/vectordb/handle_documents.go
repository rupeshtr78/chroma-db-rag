package vectordb

import (
	"chroma-db/internal/constants"
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/rs/zerolog/log"
)

// AddRecordSetToCollection adds a record set to a collection
func (c *ChromagoCollection) AddRecordSetToCollection(ctx context.Context,
	recordSet *types.RecordSet,
	docs []string,
	metadata constants.Metadata) (*chromago.Collection, error) {

	recordSet, err := addDocumentsToRecordSet(ctx, recordSet, docs, metadata)
	if err != nil {
		return nil, err
	}

	err = validateRecordSet(ctx, recordSet)
	if err != nil {
		return nil, err
	}

	coll, err := addRecordsToCollection(ctx, c, recordSet)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Added %d records to collection", len(docs))

	return coll, nil
}

// addDocumentsToRecordSet adds documents to the record set
func addDocumentsToRecordSet(ctx context.Context,
	recordSet *types.RecordSet,
	docs []string,
	metadata constants.Metadata) (*types.RecordSet, error) {

	recordSet, err := AddTextToRecordSet(ctx, recordSet, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding to record set: %v\n", err)
		return nil, err
	}
	return recordSet, nil
}

// validateRecordSet validates the record set
func validateRecordSet(ctx context.Context, recordSet *types.RecordSet) error {
	_, err := recordSet.BuildAndValidate(ctx)
	if err != nil {
		log.Debug().Msgf("Error building and validating records: %v\n", err)
		return err
	}
	return nil
}

// addRecordsToCollection adds records to the collection
func addRecordsToCollection(ctx context.Context,
	collection *ChromagoCollection,
	recordSet *types.RecordSet) (*chromago.Collection, error) {

	c, err := collection.AddRecords(ctx, recordSet)
	if err != nil {
		log.Err(err).Msg("Error adding records to collection")
		return nil, err
	}
	return c, nil
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
