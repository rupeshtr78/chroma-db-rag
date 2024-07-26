package db

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/documents"
	ollamamodel "chroma-db/internal/ollama"
	"chroma-db/internal/vectordb"
	"chroma-db/pkg/logger"
	"context"
	"strings"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/vectorstores"
)

var log = logger.Log

func QueryVectorDatabase(ctx context.Context, queryString string) (string, error) {
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

	client.SetTenant(constants.TenantName)

	// Get or create the database
	d, err := chromaclient.GetOrCreateDatabase(ctx, client, constants.Database, t.Name)
	if err != nil {
		log.Debug().Msgf("Error getting or creating database: %v\n", err)
		return "", err
	}

	log.Debug().Msgf("Database: %v\n", *d.Name)

	client.SetDatabase(constants.Database)

	// Get the ollama embedding function
	ollamaEmbedFn, err := ollamamodel.GetOllamaEmbedding(constants.OllamaUrl, constants.OllamaEmbdedModel)
	if err != nil {
		log.Debug().Msgf("Error getting ollama embedding function: %v\n", err)
		return "", err
	}

	// Create a new store
	store, err := chromaclient.CreateChromaStore(ctx,
		constants.ChromaUrl,
		constants.Namespace,
		ollamaEmbedFn,
		types.DistanceFunction(constants.DistanceFn))
	if err != nil {
		log.Debug().Msgf("Error creating store: %v\n", err)
		return "", err
	}

	// Get the list of all the available collections
	collections, err := client.ListCollections(ctx)
	if err != nil {
		log.Debug().Msgf("Error getting collections: %v\n", err)
		return "", err
	}

	// Print the list of databases

	for _, col := range collections {
		log.Debug().Msgf("Collection: %v\n", col.Name)
		log.Debug().Msgf("Database: %v\n", col.Database)
		log.Debug().Msgf("Tenant: %v\n", col.Tenant)
		log.Debug().Msgf("Embedding Function: %v\n", col.EmbeddingFunction)
	}

	// // get the documents from the pdf
	pdfDocs, err := documents.PdfToDocument(ctx, "text/Model Params.pdf")
	if err != nil {
		log.Debug().Msgf("Error getting pdf documents: %v\n", err)
		return "", err
	}

	if len(pdfDocs) == 0 {
		log.Debug().Msgf("No documents found in the pdf\n")
		return "", nil
	}

	// Add the documents to the store
	vecAddOptions := []vectorstores.Option{
		vectorstores.WithNameSpace(constants.Namespace),
	}

	// Add the documents to the store
	s, err := store.AddDocuments(ctx, pdfDocs, vecAddOptions...)
	if err != nil {
		log.Err(err).Msgf("Error adding documents to the store: %v\n", err)
		return "", err
	}

	log.Info().Msgf("Added %v documents to the store\n", s)

	vecSearchOptions := []vectorstores.Option{
		vectorstores.WithNameSpace(constants.Namespace),
		vectorstores.WithScoreThreshold(constants.ScoreThreshold),
	}

	// // Search the store
	// queryString := "what is mirostat_tau?"
	docs, err := vectordb.SearchVectorDb(ctx,
		store,
		queryString,
		3,
		constants.Namespace,
		vecSearchOptions...)

	if err != nil {
		log.Debug().Msgf("Error searching the store: %v\n", err)
		return "", err
	}

	log.Info().Msgf("Found %v documents\n", len(docs))

	var results strings.Builder
	for _, doc := range docs {
		log.Debug().Msgf("Document: %v\n", doc.PageContent)
		log.Debug().Msgf("Metadata: %v\n", doc.Metadata)
		log.Debug().Msgf("Score: %v\n", doc.Score)

		results.WriteString(doc.PageContent)

	}

	return results.String(), nil

}
