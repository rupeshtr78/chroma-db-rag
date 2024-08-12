package main

import (
	"bufio"
	"chroma-db/app/chat"
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/documenthandler"
	"chroma-db/internal/prompts"
	"chroma-db/internal/reranker"
	"chroma-db/internal/vectordbquery"
	"chroma-db/pkg/logger"
	"context"
	"flag"
	"fmt"
	"os"
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

	// Define the flags
	loadFlag := flag.Bool("load", false, "Load the data")
	queryFlag := flag.Bool("query", false, "Serve the data")

	// Parse the flags
	flag.Parse()

	if !*loadFlag && !*queryFlag {
		fmt.Println("Please specify a flag: -load or -query")
		return
	}

	// Create a new reader to read from standard input
	reader := bufio.NewReader(os.Stdin)

	var documentPath string
	// If the load flag is set, load the data
	if *loadFlag {
		fmt.Print("Enter the path to the document file: ")
		path, err := reader.ReadString('\n')
		if err != nil && err.Error() != "EOF" {
			fmt.Println("Error reading input:", err)
			return
		}
		documentPath = path[:len(path)-1] // Remove the newline character
	}

	var queryString string
	// If the query flag is set, query the data
	if *queryFlag {
		fmt.Print("Enter the query: ")
		queryStr, err := reader.ReadString('\n')
		if err != nil && err.Error() != "EOF" {
			fmt.Println("Error reading input:", err)
			return
		}
		queryString = queryStr[:len(queryStr)-1] // Remove the newline character
	}

	// Initialize the Chroma client
	client, err := chromaclient.GetChromaClientInstance(ctx, constants.ChromaUrl, constants.TenantName, constants.Database)
	if err != nil {
		log.Debug().Msgf("Error initializing Chroma: %v\n", err)
		return
	}

	errChan := make(chan error, 1)
	defer close(errChan)
	collectionChan := make(chan *chromago.Collection, 1)
	defer close(collectionChan)
	var wg sync.WaitGroup

	if documentPath != "" {
		// Load the data
		log.Debug().Msgf("Loading data from: %v", documentPath)
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

		c := <-collectionChan
		log.Debug().Msgf("Loaded data to Collection: %v", c.Name)

	}

	// queryString := "what is the difference between mirostat_tau and mirostat_eta?"
	if queryString != "" {
		// Query the collection with the query text
		vectorQuery := []string{queryString}
		vectorChan := make(chan *chromago.QueryResults, 1)
		rankChan := make(chan *reranker.HfRerankResponse, 1)
		defer close(vectorChan)

		c, err := chromaclient.GetCollection(ctx, client.Client, constants.Collection, constants.HuggingFace)
		if err != nil {
			log.Error().Msgf("Failed to get collection: %v", err)
		}

		wg.Add(1)
		go func(c context.Context, collection *chromago.Collection, query []string) {
			defer wg.Done()
			vectorResults, err := vectordbquery.QueryVectorDbWithOptions(ctx, collection, query)
			if err != nil {
				errChan <- err
				log.Error().Msgf("Failed to query vector db: %v", err)
			}
			vectorChan <- vectorResults

		}(ctx, c, vectorQuery)

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

}

func processDocument(ctx context.Context, client *chromaclient.ChromaClient, path string, docType constants.DocType, errChan chan<- error) {
	collection, err := documenthandler.VectorEmbedData(ctx, client,
		documenthandler.WithDocPath(path),
		documenthandler.WithDocType(docType))
	if err != nil {
		errChan <- err
	}
	log.Debug().Msgf("Loaded data to Collection: %v", collection.Name)
}
