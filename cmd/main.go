package main

import (
	"chroma-db/app/chat"
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/documenthandler"
	"chroma-db/internal/prompts"
	"chroma-db/internal/reranker"
	"chroma-db/internal/vectordbquery"
	"chroma-db/pkg/logger"
	"context"
	"sync"

	chromago "github.com/amikos-tech/chroma-go"
)

var log = logger.Log

func main() {
	ctx := context.Background()
	// sometimes timeout happens while model is running on remote server switching to cancel context //TODO
	// ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initialize the Chroma client
	client, err := chromaclient.GetChromaClientInstance(ctx, constants.ChromaUrl, constants.TenantName, constants.Database)
	if err != nil {
		log.Debug().Msgf("Error initializing Chroma: %v\n", err)
		return
	}

	errChan := make(chan error, 1)
	collectionChan := make(chan *chromago.Collection, 1)
	defer close(errChan)
	defer close(collectionChan)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(ctx context.Context, path string, docType constants.DocType) {
		defer wg.Done()
		collection, err := documenthandler.VectorEmbedData(ctx, client,
			documenthandler.WithDocPath(path),
			documenthandler.WithDocType(constants.TXT))
		if err != nil {
			errChan <- err
		}
		collectionChan <- collection

	}(ctx, "test/model_params.txt", constants.TXT)

	// Query the collection with the query text
	queryString := "what is the difference between mirostat_tau and mirostat_eta?"
	// queryString := "what is mirostat_eta?"
	// vectorQuery := stripStopWords(queryString) // tried this option no better results
	vectorQuery := []string{queryString}

	vectorChan := make(chan *chromago.QueryResults, 1)
	rankChan := make(chan *reranker.HfRerankResponse, 1)
	defer close(vectorChan)

	wg.Add(1)
	go func(c context.Context, query []string) {
		defer wg.Done()
		collection := <-collectionChan

		vectorResults, err := vectordbquery.QueryVectorDbWithOptions(ctx, collection, query)
		if err != nil {
			errChan <- err
			log.Error().Msgf("Failed to query vector db: %v", err)
		}
		vectorChan <- vectorResults

	}(ctx, vectorQuery)

	wg.Add(1)
	go func(c context.Context, query []string) {
		defer wg.Done()
		queryResults := <-vectorChan
		rerankResults, err := vectordbquery.RerankQueryResult(c, query, queryResults.Documents[0])
		if err != nil {
			errChan <- err
			log.Error().Msgf("Failed to rerank query results: %v", err)
		}
		rankChan <- rerankResults
	}(ctx, vectorQuery)

	// wait for all go routines to finish
	wg.Wait()

	// Get the final prompt and chat result
	var rankResult *reranker.HfRerankResponse
	select {
	case rankResult = <-rankChan:
	case <-ctx.Done():
		log.Error().Msgf("Context Timeout: %v", ctx.Err())
		return
	}

	contentString := rankResult.Text
	prompts, err := prompts.GetTemplate(constants.SystemPromptFile, queryString, contentString)
	if err != nil {
		log.Error().Msgf("Failed to get template: %v", err)

	}

	log.Debug().Msgf("Final Prompt: %v", prompts)

	chat.ChatOllama(ctx, prompts)

}
