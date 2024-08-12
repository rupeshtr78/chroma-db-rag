package vectordbquery

import (
	"chroma-db/internal/constants"
	"chroma-db/internal/reranker"
	"chroma-db/pkg/logger"
	"context"
	"fmt"
	"net/http"
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

// QueryVectorDb queries the vector database with the given query text // TODO: Add more details
func QueryVectorDb(ctx context.Context, collection *chromago.Collection, queryTexts []string) (*chromago.QueryResults, error) {
	// Query the collection
	qr, qrerr := collection.Query(ctx,
		queryTexts,
		5,
		nil,
		nil,
		nil)

	if qrerr != nil {
		log.Debug().Msgf("Error querying collection: %s \n", qrerr)
		return nil, qrerr
	}

	numResults := len(qr.Documents[0])
	if numResults == 0 {
		return nil, fmt.Errorf("no results found for query: %v", queryTexts)
	}

	log.Debug().Msgf("Query Results Length: %v\n", numResults)
	log.Debug().Msgf("Query Distance: %v\n", qr.Distances)
	log.Debug().Msgf("Query Metadata: %v\n", qr.Metadatas)

	return qr, nil
}

// QueryVectorDbWithOptions queries the vector database with the given query text and options
func QueryVectorDbWithOptions(ctx context.Context, collection *chromago.Collection, queryTexts []string) (*chromago.QueryResults, error) {
	// Query the collection
	str := strings.Builder{}
	for _, text := range queryTexts {
		str.WriteString(text)
	}

	// Embed the query text
	embedding, err := collection.EmbeddingFunction.EmbedQuery(ctx, str.String())
	if err != nil {
		log.Debug().Msgf("Error embedding query: %s \n", err)
		return nil, err
	}

	queryEmbedder := []*types.Embedding{embedding}
	// _ = queryEmbedder

	options := []types.CollectionQueryOption{
		// types.WithQueryTexts(queryTexts),
		types.WithQueryText(str.String()),
		types.WithNResults(5), // add more results need testing
		// types.WithOffset(10),
		types.WithQueryEmbeddings(queryEmbedder),
	}

	qr, qrerr := collection.QueryWithOptions(ctx, options...)
	if qrerr != nil {
		log.Debug().Msgf("Error querying collection: %s \n", qrerr)
		return nil, qrerr
	}

	numResults := len(qr.Documents[0])
	log.Debug().Msgf("Query Results Length: %v\n", numResults)
	log.Debug().Msgf("Query Distance: %v\n", qr.Distances)
	log.Debug().Msgf("Query Metadata: %v\n", qr.Metadatas)

	// docs := qr.Documents[0]
	// rerankIndex, err := RerankQueryResult(ctx, str.String(), docs)
	// if err != nil {
	// 	log.Error().Msgf("Error reranking query results: %v\n", err)
	// }
	// index := rerankIndex.Index
	// TODO may be add reranking logic here
	// assuming smaller distance is better pick the first result qr.Documents[0][0]
	// adding the second result qr.Documents[0][1] better results
	// concatenate the results qr.Documents[0][0] and qr.Documents[0][1]
	// For specific query results with lowest distance is better qr.Documents[0][0]
	// When asked general questions, trying concatenating the results for now
	// queryResults := qr.Documents[0][1] + qr.Documents[0][0]
	// queryResults := qr.Documents[0][index]
	return qr, nil
}

// RerankQueryResult reranks the query results using the HuggingFace reranker
// TODO: Use GRPC https://github.com/huggingface/text-embeddings-inference?tab=readme-ov-file#grpc
func RerankQueryResult(ctx context.Context, queryTexts []string, queryResults []string) (*reranker.HfRerankResponse, error) {

	queryString := strings.Builder{}
	for _, text := range queryTexts {
		queryString.WriteString(text)
	}
	request := &reranker.HfRerankRequest{
		Query:       queryString.String(),
		Texts:       queryResults,
		RawScores:   false,
		ReturnTexts: true,
	}

	client := &reranker.HfRerankClient{
		Client:  &http.Client{},
		BaseURL: constants.HuggingFaceRerankUrl,
		Model:   constants.HuggingFaceRerankModel,
	}

	res, err := client.CreateRerankingRequest(ctx, request)
	if err != nil {
		log.Error().Msgf("Error reranking query results: %v\n", err)
	}

	log.Info().Msgf("Reranked Results: %v\n", res)
	// For now return the first result
	firstResult := (*res)[0]
	return &firstResult, nil
}
