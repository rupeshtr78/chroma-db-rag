package ollamarag

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/pkg/logger"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

var log = logger.Log

// AddRecordSetToCollection adds the record set to the collection.
func AddRecordSetToCollection(ctx context.Context,
	collection *chromago.Collection,
	recordSet *types.RecordSet,
	docs []string,
	metadata constants.Metadata) (*chromago.Collection, error) {
	// Add the documents to the record set
	recordSet, err := chromaclient.AddTextToRecordSet(ctx, collection, recordSet, docs, metadata)
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
