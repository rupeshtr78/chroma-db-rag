package chromaclient

import (
	"chroma-db/internal/constants"
	"chroma-db/internal/embedders"
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/types"
)

// CreateCollection creates a new collection with the given name, embedding function and distance function.
func CreateCollection(ctx context.Context, chromaClient *ChromaClient, hfef types.EmbeddingFunction) (*chromago.Collection, error) {

	// delete the collection if it exists
	err := DeleteCollectionIfExists(ctx, constants.Collection, chromaClient.Client, hfef)
	if err != nil {
		log.Debug().Msgf("Error deleting collection: %v\n", err)
		return nil, err
	}

	// Create a new collection with the given name client tenant and database
	collection, err := CreateNewCollection(ctx, chromaClient.Client,
		constants.Collection,
		hfef,
		constants.DistanceFn)
	if err != nil {
		log.Debug().Msgf("Error getting or creating collection: %v\n", constants.Collection)
		return nil, err
	}

	return collection, nil
}

// CreateNewCollection creates a new **chromago.Collection** if it does not exist
func CreateNewCollection(ctx context.Context,
	client *chromago.Client,
	collectionName string,
	embeddingFunction types.EmbeddingFunction,
	distanceFn types.DistanceFunction) (*chromago.Collection, error) {

	// Create a new collection with options
	newCollection, err := client.NewCollection(
		ctx,
		collection.WithName(collectionName),
		collection.WithCreateIfNotExist(true),
		collection.WithEmbeddingFunction(embeddingFunction),
		collection.WithHNSWDistanceFunction(distanceFn),
		collection.WithTenant(constants.TenantName),
		collection.WithDatabase(constants.Database),
	)
	if err != nil {
		log.Err(err).Msg("error creating collection")
		return nil, err
	}

	log.Debug().Msgf("Collection %v created\n", collectionName)

	newCollection.Tenant = constants.TenantName
	newCollection.Database = constants.Database

	return newCollection, nil
}

func DeleteCollectionIfExists(ctx context.Context, collectionName string, client *chromago.Client, embeddingFunction types.EmbeddingFunction) error {

	// List all collections Check if the collection already exist
	collections, err := client.ListCollections(ctx)
	if err != nil {
		log.Debug().Msgf("Error listing collections: %v\n", err)
		return err
	}
	for _, c := range collections {
		if c.Name == collectionName {
			// Collection already exists, Delete the collection
			collection, err := client.DeleteCollection(ctx, collectionName)
			if err != nil {
				log.Err(err).Msgf("Error deleting collection: %s \n", collectionName)
				return err
			}
			log.Debug().Msgf("Collection %v deleted\n", collection.Name)
		}
	}

	return nil
}

func GetCollection(ctx context.Context,
	client *chromago.Client,
	collectionName string) (*chromago.Collection, error) {

	// List all collections Check if the collection already exist
	collections, err := client.ListCollections(ctx)
	if err != nil {
		log.Debug().Msgf("Error listing collections: %v\n", err)
		return nil, err
	}

	for _, c := range collections {
		if c.Name == collectionName {
			return c, nil
		}
	}

	return nil, fmt.Errorf("collection %v not found", collectionName)

}

func GetCollectionFromDb(ctx context.Context, chromaClient *chromago.Client, embbedder constants.Embedder, embeddingModel string) (*chromago.Collection, error) {
	// Get Embedding either HuggingFace or Ollama
	em := embedders.NewEmbeddingManager(embbedder, constants.HuggingFaceTeiUrl, embeddingModel)

	hfef, err := em.GetEmbeddingFunction()
	if err != nil {
		log.Debug().Msgf("Error getting hugging face embedding function: %v\n", err)
		return nil, err
	}

	// Create a new collection with the given name client tenant and database
	collection, err := CreateNewCollection(ctx, chromaClient,
		constants.Collection,
		hfef,
		constants.DistanceFn)
	if err != nil {
		log.Debug().Msgf("Error getting or creating collection: %v\n", constants.Collection)
		return nil, err
	}

	return collection, nil

}
