package vectordb

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/embedders"
	"chroma-db/pkg/logger"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

var log = logger.Log

// InitializeClient initializes the Chroma client and sets the tenant and database.
// Creates a new collection with the given name, embedding function and distance function.
// Creates a new record set.
// Returns the collection and record set. // TODO refactor doing too much
func InitializeChroma(ctx context.Context, chromaUrl string, tenantName string, databaseName string) (*chromago.Client, error) {
	// Initialize the chroma client
	client, err := chromaclient.GetChromaClient(ctx, constants.ChromaUrl)
	if err != nil {
		log.Debug().Msgf("Error getting chroma client: %v\n", err)
		return nil, err
	}

	// // Get or create the tenant
	_, err = chromaclient.GetOrCreateTenant(ctx, client, constants.TenantName)
	if err != nil {
		log.Debug().Msgf("Error getting or creating tenant: %v\n", err)
		return nil, err
	}

	// Set the tenant for the client
	client.SetTenant(constants.TenantName)

	// // Get or create the database
	_, err = chromaclient.GetOrCreateDatabase(ctx, client, constants.Database, &constants.TenantName)
	if err != nil {
		log.Debug().Msgf("Error getting or creating database: %v\n", err)
		return nil, err
	}

	// Set the database for the client
	client.SetDatabase(constants.Database)

	// client.SetDatabase(constants.Database)
	log.Debug().Msgf("Client Tenant: %v\n", client.Tenant)
	log.Debug().Msgf("Client Database: %v\n", client.Database)

	return client, nil

}

func GetCollectionRecordSet(ctx context.Context, client *chromago.Client, embbedder constants.Embedder, embeddingModel string) (*chromago.Collection, *types.RecordSet, error) {
	// Get Embedding either HuggingFace or Ollama
	em := embedders.NewEmbeddingManager(embbedder, constants.HuggingFaceTeiUrl, embeddingModel)

	hfef, err := em.GetEmbeddingFunction()
	if err != nil {
		log.Debug().Msgf("Error getting hugging face embedding function: %v\n", err)
		return nil, nil, err
	}

	// delete the collection if it exists
	err = chromaclient.DeleteCollectionIfExists(ctx, constants.Collection, client, hfef)
	if err != nil {
		log.Debug().Msgf("Error deleting collection: %v\n", err)
		return nil, nil, err
	}

	// Create a new collection with the given name client tenant and database
	collection, err := chromaclient.GetOrCreateCollection(ctx, client,
		constants.Collection,
		hfef,
		constants.DistanceFn)
	if err != nil {
		log.Debug().Msgf("Error getting or creating collection: %v\n", constants.Collection)
		return nil, nil, err
	}

	// Create a new record set
	recordSet, err := chromaclient.CreateRecordSet(hfef)
	if err != nil {
		log.Debug().Msgf("Error creating record set: %v\n", err)
		return nil, nil, err
	}

	return collection, recordSet, nil
}
