// POC to rerank results based on distance
package reranker

import (
	"context"
	"errors"
	"fmt"
	"sort"

	chromago "github.com/amikos-tech/chroma-go"
)

type SimpleReranker struct{}

func (sr *SimpleReranker) Rerank(ctx context.Context, query string, queryResults *chromago.QueryResults) ([]*RankedResult, error) {
	if len(queryResults.Documents) == 0 {
		return nil, errors.New("no results to rerank")
	}

	// rank based on distances
	queryDistances := queryResults.Distances[0]
	// fetch index of queryResults.Distances[0] between 0.35 and 0.40
	for idx, dist := range queryDistances {
		if dist >= 0.35 && dist <= 0.40 {
			fmt.Printf("idx: %v, dist: %v\n", idx, dist)
		}
	}

	rankedResults := make([]*RankedResult, len(queryResults.Documents))
	for idx, result := range queryResults.Documents {
		rankedResults[idx] = &RankedResult{
			ID:     idx,
			String: result[idx],
			Rank:   queryDistances[idx], // rank based on distance
		}
	}

	// Sort by rank (ascending, assuming smaller distance is better)
	sort.Slice(rankedResults, func(i, j int) bool {
		return rankedResults[i].Rank < rankedResults[j].Rank
	})

	return rankedResults, nil
}

func (sr *SimpleReranker) RerankResults(ctx context.Context, queryResults *chromago.QueryResults) (RerankedChromaResults, error) {
	// Example logic to rerank results
	if queryResults == nil || len(queryResults.Documents) == 0 {
		return RerankedChromaResults{}, errors.New("no query results to rerank")
	}

	reranked := RerankedChromaResults{
		QueryResults: *queryResults,
		Ranks:        make([][]float32, len(queryResults.Documents)),
	}

	for idx, result := range queryResults.Documents {
		reranked.Ranks[idx] = []float32{float32(len(result))}
	}

	return reranked, nil
}
