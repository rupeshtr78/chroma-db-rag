package chromaclient

import (
	"context"
	"log"

	chromago "github.com/amikos-tech/chroma-go"
	openapi "github.com/amikos-tech/chroma-go/swagger"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

func GetChromaClient(ctx context.Context, url string) (*chromago.Client, error) {
	// create the client connection and confirm that we can access the server with it
	chromaClient, err := chromago.NewClient(url)
	if err != nil {
		return nil, err
	}

	if _, errHb := chromaClient.Heartbeat(ctx); errHb != nil {
		return nil, errHb
	}

	return chromaClient, err
}

func CreateChromaStore(ctx context.Context,
	chromaUrl string,
	nameSpace string,
	embedder embeddings.Embedder,
	distanceFunction types.DistanceFunction) (*chroma.Store, error) {

	store, err := chroma.New(
		chroma.WithChromaURL(chromaUrl),
		chroma.WithNameSpace(nameSpace),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(distanceFunction), // default is cosine l2 ip
	)
	if err != nil {
		return nil, err
	}

	return &store, nil

}

func GetOrCreateTenant(ctx context.Context, client *chromago.Client, tenantName string) (*openapi.Tenant, error) {

	if t, err := client.GetTenant(ctx, tenantName); err == nil {
		log.Default().Printf("Tenant %v already exists\n", tenantName)
		return t, nil
	}

	t, err := client.CreateTenant(ctx, tenantName)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}
	return t, nil
}

func GetOrCreateDatabase(ctx context.Context, client *chromago.Client, dbName string, tenantName *string) (*openapi.Database, error) {

	if d, err := client.GetDatabase(ctx, dbName, tenantName); err == nil {
		log.Default().Printf("Database %v already exists\n", dbName)
		return d, nil
	}

	d, err := client.CreateDatabase(ctx, dbName, tenantName)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}
	return d, nil
}

// GetOrCreateCollection creates a new **chromago.Collection** if it does not exist
func GetOrCreateCollection(ctx context.Context,
	client *chromago.Client,
	collectionName string,
	embeddingFunction types.EmbeddingFunction,
	distanceFn types.DistanceFunction) (*chromago.Collection, error) {

	if c, err := client.GetCollection(ctx,
		collectionName,
		embeddingFunction); err == nil {
		log.Default().Printf("Collection %v already exists\n", collectionName)
		return c, nil
	}

	c, err := client.CreateCollection(ctx,
		collectionName,
		nil,
		true,
		embeddingFunction,
		distanceFn)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}
	return c, nil
}
