package ollamarag

import (
	"chroma-db/cmd/vectordb"
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/documents"
	"chroma-db/pkg/logger"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

var log = logger.Log

// LoadDocuments loads and returns documents and metadata from a given path.
func LoadTextDocuments(ctx context.Context, docPath string) ([]string, constants.Metadata, error) {
	docs, metadata, err := documents.TextLoaderV2(ctx, docPath)
	if err != nil {
		log.Debug().Msgf("Error loading text: %v\n", err)
		return nil, nil, err
	}

	return docs, metadata, nil
}

func LoadPdfDocuments(ctx context.Context, docPath string) ([]string, constants.Metadata, error) {
	docs, metadata, err := documents.ParsePDF(docPath)
	if err != nil {
		log.Debug().Msgf("Error parsing PDF: %v\n", err)
		return nil, nil, err
	}

	return docs, metadata, nil
}

// AddRecordSetToCollection adds the record set to the collection.
func AddRecordSetToCollection(ctx context.Context,
	collection *chromago.Collection,
	recordSet *types.RecordSet, docs []string, metadata constants.Metadata) (*chromago.Collection, error) {
	// Add the documents to the record set
	recordSet, err := chromaclient.AddTextToRecordSet(ctx, collection, recordSet, docs, metadata)
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

	log.Debug().Msgf("Added %d records to collection: %s\n", len(docs), collection.Name)

	return collection, nil
}

func RunOllamaRag(ctx context.Context, chromaUrl string, tenantName string, databaseName string, docPath string) (*chromago.Collection, error) {
	// Initialize the Chroma client
	collection, recordSet, err := vectordb.InitializeChroma(ctx, chromaUrl, tenantName, databaseName)
	if err != nil {
		log.Debug().Msgf("Error initializing Chroma: %v\n", err)
		return nil, err
	}

	// Load the documents
	docs, metadata, err := LoadTextDocuments(ctx, docPath)
	if err != nil {
		log.Debug().Msgf("Error loading documents: %v\n", err)
		return nil, err
	}

	// Add the record set to the collection
	collection, err = AddRecordSetToCollection(ctx, collection, recordSet, docs, metadata)
	if err != nil {
		log.Debug().Msgf("Error adding record set to collection: %v\n", err)
		return nil, err
	}

	// Count the number of documents in the collection
	countDocs, qrerr := collection.Count(ctx)
	if qrerr != nil {
		log.Debug().Msgf("Error counting documents: %s \n", qrerr)
	}

	log.Debug().Msgf("Number of documents in collection: %d\n", countDocs)

	return collection, nil
}
