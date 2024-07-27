package chromaclient

import (
	"context"

	chromago "github.com/amikos-tech/chroma-go"
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

	c, err := client.CreateCollection(ctx,
		collectionName,
		Metadata{},
		true,
		embeddingFunction,
		distanceFn)
	if err != nil {
		log.Err(err).Msgf("Failed to create collection %v\n", collectionName)
		return nil, err
	}
	return c, nil
}
