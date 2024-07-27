package db

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/documents"
	ollamamodel "chroma-db/internal/ollama"
	"context"
)

func LoadDataToVectorDB(ctx context.Context, docPath string) (string, error) {
	// Get the chroma client
	client, err := chromaclient.GetChromaClient(ctx, constants.ChromaUrl)
	if err != nil {
		log.Debug().Msgf("Error getting chroma client: %v\n", err)
		return "", err
	}

	// Get or create the tenant
	t, err := chromaclient.GetOrCreateTenant(ctx, client, constants.TenantName)
	if err != nil {
		log.Debug().Msgf("Error getting or creating tenant: %v\n", err)
		return "", err
	}

	client.SetTenant(*t.Name)

	// Get or create the database
	_, err = chromaclient.GetOrCreateDatabase(ctx, client, constants.Database, t.Name)
	if err != nil {
		log.Debug().Msgf("Error getting or creating database: %v\n", err)
		return "", err
	}

	log.Debug().Msgf("Database: %v\n", constants.Database)

	client.SetDatabase(constants.Database)

	ollamaEmbedFn, err := ollamamodel.GetOllamaEmbeddingFn(constants.OllamaUrl, constants.OllamaEmbdedModel)
	if err != nil {
		log.Debug().Msgf("Error getting ollama embedding function: %v\n", err)
		return "", err
	}

	collection, err := chromaclient.GetOrCreateCollection(ctx, client,
		constants.Collection,
		ollamaEmbedFn,
		constants.DistanceFn)
	if err != nil {
		log.Debug().Msgf("Error getting or creating collection: %v\n", constants.Collection)
		return "", err
	}

	recordSet, err := chromaclient.CreateRecordSet(ollamaEmbedFn)
	if err != nil {
		log.Debug().Msgf("Error creating record set: %v\n", err)
		return "", err
	}

	docs, metadata, err := documents.ParsePDF(docPath)
	if err != nil {
		log.Debug().Msgf("Error parsing PDF: %v\n", err)
		return "", err
	}

	recordSet, err = chromaclient.AddToRecordSet(ctx, collection, recordSet, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding to record set: %v\n", err)
		return "", err
	}

	// Add the records to the collection
	collection, err = collection.AddRecords(ctx, recordSet)
	if err != nil {
		log.Err(err).Msgf("Error adding records to collection: %s\n", collection.Name)
		return "", err
	}

	return "", nil

}
