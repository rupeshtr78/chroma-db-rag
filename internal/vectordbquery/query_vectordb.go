package vectordbquery

import (
	"chroma-db/pkg/logger"
	"context"
	"fmt"
	"strings"

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

func QueryVectorDb(ctx context.Context, collection *chromago.Collection, queryTexts []string) (string, error) {
	// Query the collection
	qr, qrerr := collection.Query(ctx,
		queryTexts,
		5,
		nil,
		nil,
		nil)

	if qrerr != nil {
		log.Debug().Msgf("Error querying collection: %s \n", qrerr)
		return "", qrerr
	}
	fmt.Printf("qr: %v\n", qr.Documents[0][0]) //this should result in the document about dogs
	log.Info().Msgf("Query Distance: %v\n", qr.Distances)
	log.Info().Msgf("Query Metadata: %v\n", qr.Metadatas)

	queryResults := qr.Documents[0][0]

	return queryResults, nil
}

func QueryVectorDbWithOptions(ctx context.Context, collection *chromago.Collection, queryTexts []string) (string, error) {
	// // TODO Remove added for Poc for embedding query
	// embed := &types.Embedding{
	// 	ArrayOfFloat32: &[]float32{0.3, 0.4, 0.6},
	// 	ArrayOfInt32:   &[]int32{0, 2, 4},
	// }

	// queryTexts = []string{"what is the difference between mirostat_tau and mirostat_eta?"}

	str := strings.Builder{}
	for _, text := range queryTexts {
		str.WriteString(text)
	}

	embedding, err := collection.EmbeddingFunction.EmbedQuery(ctx, str.String())
	if err != nil {
		log.Debug().Msgf("Error embedding query: %s \n", err)
		return "", err
	}

	queryEmbedder := []*types.Embedding{embedding}
	// _ = queryEmbedder

	options := []types.CollectionQueryOption{
		// types.WithQueryTexts(queryTexts),
		types.WithQueryText(str.String()),
		types.WithNResults(2), // add more results need testing
		// types.WithOffset(10),
		types.WithQueryEmbeddings(queryEmbedder),
	}

	qr, qrerr := collection.QueryWithOptions(ctx, options...)
	if qrerr != nil {
		log.Debug().Msgf("Error querying collection: %s \n", qrerr)
		return "", qrerr
	}

	numResults := len(qr.Documents[0])
	log.Debug().Msgf("Query Results Length: %v\n", numResults)

	fmt.Printf("qr: %v\n", qr.Documents[0][1])
	log.Info().Msgf("Query Distance: %v\n", qr.Distances)
	log.Info().Msgf("Query Metadata: %v\n", qr.Metadatas)

	// TODO may be add reranking logic here
	// assuming smaller distance is better pick the first result qr.Documents[0][0]
	// trying to get the second result qr.Documents[0][1] better results
	// concatenate the results qr.Documents[0][0] and qr.Documents[0][1]
	// For specific query results with lowest distance is better qr.Documents[0][0]
	// When asked general questions, trying concatenating the results for now
	queryResults := qr.Documents[0][1] + qr.Documents[0][0]
	return queryResults, nil
}

// TODO: Implement this function
func RerankQueryResult() {
	// queryRerank := &reranker.SimpleReranker{}

	// rankedResults, err := queryRerank.Rerank(ctx, "", qr)
	// if err != nil {
	// 	fmt.Println("Error in reranking:", err)
	// 	return "", err
	// }

	// rankedResults, err := queryRerank.RerankResults(ctx, qr)
	// if err != nil {
	// 	log.Debug().Msgf("Error reranking query results: %s \n", err)
	// 	return "", err
	// }
}
