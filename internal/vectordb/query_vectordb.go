package vectordb

import (
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/rs/zerolog/log"
)

// EmbedQuery embeds the query text and returns the embedding or an error //TODO: fix payload too large error
func EmbedQuery(ctx context.Context, embeddingFunc EmbeddingFunc, query []string) ([]*types.Embedding, error) {
	embedding, err := embeddingFunc.EmbedDocuments(ctx, query)
	if err != nil {
		log.Debug().Msgf("Error embedding query: %s \n", err)
		return nil, err
	}
	return embedding, nil
}

// QueryVectorDbWithOptions queries the vector database with the given query text and options
func QueryVectorDbWithOptions(ctx context.Context, collection *ChromagoCollection, queryTexts []string) (*chromago.QueryResults, error) {
	// Query the collection
	queryEmbeddings, err := EmbedQuery(ctx, collection.EmbeddingFunction(), queryTexts)
	if err != nil {
		return nil, err
	}

	// log.Debug().Msgf("Query Embeddings: %v\n", queryEmbeddings)

	options := []types.CollectionQueryOption{
		types.WithQueryTexts(queryTexts),
		// types.WithQueryText(str.String()),
		types.WithNResults(5), // add more results
		// types.WithOffset(10),
		types.WithQueryEmbeddings(queryEmbeddings), // error payload too large
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
