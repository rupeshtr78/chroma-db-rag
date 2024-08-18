package vectordb

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/embedders"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/rs/zerolog/log"
)

// CreateCollectionAndRecordSet creates a new collection with the given name, embedding function and distance function.
func CreateCollectionAndRecordSet(ctx context.Context, chromaClient *chromaclient.ChromaClient, embbedder constants.Embedder, embeddingModel string) (*chromago.Collection, *types.RecordSet, error) {
	// Get Embedding either HuggingFace or Ollama
	em := embedders.NewEmbeddingManager(embbedder, constants.HuggingFaceTeiUrl, embeddingModel)

	hfef, err := em.GetEmbeddingFunction()
	if err != nil {
		log.Debug().Msgf("Error getting hugging face embedding function: %v\n", err)
		return nil, nil, err
	}

	// delete the collection if it exists
	err = chromaclient.DeleteCollectionIfExists(ctx, constants.Collection, chromaClient.Client, hfef)
	if err != nil {
		log.Debug().Msgf("Error deleting collection: %v\n", err)
		return nil, nil, err
	}

	// Create a new collection with the given name client tenant and database
	collection, err := chromaclient.GetOrCreateCollection(ctx, chromaClient.Client,
		constants.Collection,
		hfef,
		constants.DistanceFn)
	if err != nil {
		log.Debug().Msgf("Error getting or creating collection: %v\n", constants.Collection)
		return nil, nil, err
	}

	// Create a new record set
	recordSet, err := CreateRecordSet(hfef)
	if err != nil {
		log.Debug().Msgf("Error creating record set: %v\n", err)
		return nil, nil, err
	}

	return collection, recordSet, nil
}

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
