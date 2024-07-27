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

	if c, err := client.GetCollection(ctx,
		collectionName,
		embeddingFunction); err == nil {
		log.Debug().Msgf("Collection %v already exists\n", collectionName)
		return c, nil
	}

	// Create a new collection with options
	newCollection, err := client.NewCollection(
		ctx,
		collection.WithName(collectionName),
		collection.WithEmbeddingFunction(embeddingFunction),
		collection.WithHNSWDistanceFunction(distanceFn),
	)
	if err != nil {
		log.Err(err).Msg("error creating collection")
		return nil, err
	}
	return newCollection, nil
}

func DeleteCollection(ctx context.Context, collectionName string, client *chromago.Client) error {
	// Check if the collection already exists
	_, err := client.GetCollection(ctx, collectionName, nil)
	if err != nil {
		log.Err(err).Msgf("Error getting collection: %s \n", collectionName)
		return err
	}

	// Collection already exists, Delete the collection
	_, err = client.DeleteCollection(ctx, collectionName)
	if err != nil {
		log.Err(err).Msgf("Error deleting collection: %s \n", collectionName)
		return err
	}
	return nil
}
