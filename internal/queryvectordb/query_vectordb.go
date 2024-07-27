package queryvectordb

import (
	"chroma-db/pkg/logger"
	"context"
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

type CollectionQuery struct {
	QueryTexts    []string
	Where         map[string]interface{}
	WhereDocument map[string]interface{}
	NResults      int32
	Offset        int32
	Limit         int32
	Ids           []string
}

var log = logger.Log

func QueryVectorDb(ctx context.Context, collection *chromago.Collection, queryTexts []string) error {
	// Query the collection
	qr, qrerr := collection.Query(ctx,
		queryTexts,
		5,
		nil,
		nil,
		nil)

	if qrerr != nil {
		log.Debug().Msgf("Error querying collection: %s \n", qrerr)
		return qrerr
	}
	fmt.Printf("qr: %v\n", qr.Documents[0][0]) //this should result in the document about dogs
	log.Info().Msgf("Query Distance: %v\n", qr.Distances)
	log.Info().Msgf("Query Metadata: %v\n", qr.Metadatas)
	return nil
}

func QueryVectorDbWithOptions(ctx context.Context, collection *chromago.Collection, queryTexts []string) error {
	// Query the collection
	// Query the collection using QueryWithOptions
	options := []types.CollectionQueryOption{
		types.WithQueryTexts(queryTexts),
		types.WithNResults(5),
		types.WithOffset(1),
	}

	qr, qrerr := collection.QueryWithOptions(ctx, options...)
	if qrerr != nil {
		log.Debug().Msgf("Error querying collection: %s \n", qrerr)
		return qrerr
	}

	fmt.Printf("qr: %v\n", qr.Documents[0][0]) // this should result in the document about dogs
	log.Info().Msgf("Query Distance: %v\n", qr.Distances)
	log.Info().Msgf("Query Metadata: %v\n", qr.Metadatas)
	return nil
}
