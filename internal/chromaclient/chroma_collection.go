package chromaclient

import (
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/types"
)

// GetOrCreateCollection creates a new **chromago.Collection** if it does not exist
func GetOrCreateCollection(ctx context.Context,
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
		collection.WithTenant(client.Tenant),
		collection.WithDatabase(client.Database),
	)
	if err != nil {
		log.Err(err).Msg("error creating collection")
		return nil, err
	}

	log.Debug().Msgf("Collection %v created\n", collectionName)

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
