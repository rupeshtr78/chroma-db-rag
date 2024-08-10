package main

import (
	"chroma-db/app/chat"
	"chroma-db/app/ollamarag"
	"chroma-db/internal/constants"
	"chroma-db/internal/prompts"
	"chroma-db/internal/vectordbquery"
	"chroma-db/pkg/logger"
	"context"
	"fmt"
	"strings"
	"sync"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/bbalet/stopwords"
)

var log = logger.Log

func main() {
	ctx := context.Background()
	// sometimes timeout happens while model is running on remote server switching to cancel context //TODO
	// ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)
	collectionChan := make(chan *chromago.Collection, 1)
	defer close(errChan)
	defer close(collectionChan)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(ctx context.Context, path string, docType constants.DocType) {
		defer wg.Done()
		collection, err := ollamarag.RunOllamaRagV2(ctx,
			ollamarag.WithDocPath(path),
			ollamarag.WithDocType(constants.TXT))
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

	vectorChan := make(chan string, 1)
	defer close(vectorChan)

	select {
	case err := <-errChan:
		log.Error().Msgf("Failed to run OllamaRag: %v", err)
		return
	case <-ctx.Done():
		log.Error().Msg("Timeout")
		return
	case collection := <-collectionChan:
		wg.Add(1)
		go func(c context.Context, collection *chromago.Collection, query []string) {
			defer wg.Done()
			vectorResults, err := vectordbquery.QueryVectorDbWithOptions(ctx, collection, query)
			if err != nil {
				errChan <- err
				log.Error().Msgf("Failed to query vector db: %v", err)
			}
			vectorChan <- vectorResults
		}(ctx, collection, vectorQuery)

	}

	// wait for all go routines to finish
	wg.Wait()

	// Get the vector results
	vectorResults := <-vectorChan
	prompts, err := prompts.GetTemplate(constants.SystemPromptFile, queryString, vectorResults)
	if err != nil {
		log.Error().Msgf("Failed to get template: %v", err)

	}

	log.Info().Msgf("Final Prompt: %v", prompts)

	chat.ChatOllama(ctx, prompts)

}

// stripStopWords removes the stop words from the text and returns a slice of words
func stripStopWords(text string) []string {
	langCode := "en"

	// remove stopwords
	cleanContent := stopwords.CleanString(text, langCode, true)
	fmt.Println(cleanContent)

	// convert to slice of words
	result := make([]string, 0)
	// split the text into words and trim the spaces
	for _, word := range strings.Split(cleanContent, " ") {
		trimmedWord := strings.TrimSpace(word)
		// remove extra spaces
		if len(trimmedWord) > 0 {
			result = append(result, trimmedWord)
		}
	}

	log.Debug().Msgf("Vector Query Strings: %v", result)

	return result
}
