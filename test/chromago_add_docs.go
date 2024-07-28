package test

import (
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/rs/zerolog/log"
)

type Metadata map[string]interface{}

// AddDocsToVectorDb adds documents to the vector store
func AddDocsToCollection(ctx context.Context,
	collection *chromago.Collection,
	metadata Metadata,
	documents []string,
	ids []string) error {

	metadataSlice := []map[string]interface{}{metadata}

	// Add documents to the vector store.
	_, err := collection.Add(context.TODO(),
		nil, metadataSlice, documents, ids)
	if err != nil {
		log.Debug().Msgf("error adding documents: %v\n", err)
		return err
	}

	return nil
}
