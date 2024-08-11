package chromaclient

import (
	"chroma-db/internal/constants"
	"chroma-db/internal/embedders"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

// CreateCollectionAndRecordSet creates a new collection with the given name, embedding function and distance function.
func CreateCollectionAndRecordSet(ctx context.Context, chromaClient *ChromaClient, embbedder constants.Embedder, embeddingModel string) (*chromago.Collection, *types.RecordSet, error) {
	// Get Embedding either HuggingFace or Ollama
	em := embedders.NewEmbeddingManager(embbedder, constants.HuggingFaceTeiUrl, embeddingModel)

	hfef, err := em.GetEmbeddingFunction()
	if err != nil {
		log.Debug().Msgf("Error getting hugging face embedding function: %v\n", err)
		return nil, nil, err
	}

	// delete the collection if it exists
	err = DeleteCollectionIfExists(ctx, constants.Collection, chromaClient.Client, hfef)
	if err != nil {
		log.Debug().Msgf("Error deleting collection: %v\n", err)
		return nil, nil, err
	}

	// Create a new collection with the given name client tenant and database
	collection, err := GetOrCreateCollection(ctx, chromaClient.Client,
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
