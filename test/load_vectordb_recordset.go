package test

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/documenthandler"
	ollamamodel "chroma-db/internal/ollama"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/rs/zerolog/log"
)

func LoadDataToVectorDB(ctx context.Context, docPath string) (*chromago.Collection, error) {
	// Initialize the chroma client
	client, err := chromaclient.GetChromaClient(ctx, constants.ChromaUrl)
	if err != nil {
		log.Debug().Msgf("Error getting chroma client: %v\n", err)
		return nil, err
	}

	// Get or create the tenant
	t, err := chromaclient.GetOrCreateTenant(ctx, client, constants.TenantName)
	if err != nil {
		log.Debug().Msgf("Error getting or creating tenant: %v\n", err)
		return nil, err
	}

	// Set the tenant for the client
	client.SetTenant(*t.Name)

	// Get or create the database
	_, err = chromaclient.GetOrCreateDatabase(ctx, client, constants.Database, t.Name)
	if err != nil {
		log.Debug().Msgf("Error getting or creating database: %v\n", err)
		return nil, err
	}

	client.SetDatabase(constants.Database)
	log.Debug().Msgf("Client Tenant: %v\n", client.Tenant)
	log.Debug().Msgf("Client Database: %v\n", client.Database)

	// Get the ollama embedding function
	ollamaEmbedFn, err := ollamamodel.GetOllamaEmbeddingFn(constants.OllamaUrl, constants.OllamaEmbdedModel)
	if err != nil {
		log.Debug().Msgf("Error getting ollama embedding function: %v\n", err)
		return nil, err
	}

	// delete the collection if it exists
	err = chromaclient.DeleteCollectionIfExists(ctx, constants.Collection, client, ollamaEmbedFn)
	if err != nil {
		log.Debug().Msgf("Error deleting collection: %v\n", err)
		return nil, err
	}

	// Create a new collection with the given name client tenant and database
	collection, err := chromaclient.GetOrCreateCollection(ctx, client,
		constants.Collection,
		ollamaEmbedFn,
		constants.DistanceFn)
	if err != nil {
		log.Debug().Msgf("Error getting or creating collection: %v\n", constants.Collection)
		return nil, err
	}

	// Create a new record set
	recordSet, err := chromaclient.CreateRecordSet(ollamaEmbedFn)
	if err != nil {
		log.Debug().Msgf("Error creating record set: %v\n", err)
		return nil, err
	}

	// Load text from a file
	docType := constants.TXT
	docLoader := documenthandler.NewDocumentLoader(docType)

	docs, metadata, err := docLoader.LoadDocument(ctx, docPath)

	// docs, metadata, err := documents.TextLoaderV2(ctx, docPath)
	if err != nil {
		log.Debug().Msgf("Error loading text: %v\n", err)
		return nil, err
	}

	// for i, doc := range docs {
	// 	key := fmt.Sprintf("%d", i+1)
	// 	log.Debug().Msgf("Document: %v\n", doc)
	// 	log.Debug().Msgf("Metadata: %v\n", metadata[key])

	// }

	// Add the documents to the record set
	recordSet, err = chromaclient.AddTextToRecordSet(ctx, collection, recordSet, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding to record set: %v\n", err)
		return nil, err
	}

	// Build and validate the record set
	_, err = recordSet.BuildAndValidate(ctx)
	if err != nil {
		log.Debug().Msgf("Error building and validating records: %v\n", err)
		return nil, err
	}

	// Add the records to the collection
	collection, err = collection.AddRecords(ctx, recordSet)
	if err != nil {
		log.Err(err).Msgf("Error adding records to collection: %s\n", collection.Name)
		return nil, err
	}

	log.Info().Msgf("Added %d records to collection: %s\n", len(docs), collection.Name)

	// Count the number of documents in the collection
	countDocs, qrerr := collection.Count(context.TODO())
	if qrerr != nil {
		log.Debug().Msgf("Error counting documents: %s \n", qrerr)
	}

	log.Info().Msgf("Number of documents in collection: %d\n", countDocs)

	return collection, nil

}
