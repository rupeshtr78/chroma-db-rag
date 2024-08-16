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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Define the flags for the application
	loadFlag := flag.Bool("load", false, "Load and embed the data in vectordb")
	queryFlag := flag.Bool("query", false, "Query the embedded data and rerank the results")
	// grpcFlag := flag.Bool("grpc", false, "Query the embedded data and rerank the results using gRPC")

	// Parse the flags
	flag.Parse()

	if !*loadFlag && !*queryFlag {
		fmt.Println("Please specify a flag: -load or -query")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var documentPath string
	if *loadFlag {
		documentPath = readInput(reader, "Enter the path to the document file: ")
	}

	var queryString string
	if *queryFlag {
		queryString = readInput(reader, "Enter the query: ")
	}

	// Initialize the Chroma client
	client, err := chromaclient.GetChromaClientInstance(ctx, constants.ChromaUrl, constants.TenantName, constants.Database)
	if err != nil {
		log.Debug().Msgf("Error initializing Chroma: %v\n", err)
		return
	}

	// Initialize the ReRank client
	reRankClient, err := reranker.NewReRanker(ctx, constants.GRPC)
	if err != nil {
		log.Error().Msgf("Error initializing ReRank Client: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	collectionChan := make(chan *chromago.Collection, 1)
	defer close(errChan)
	defer close(collectionChan)

	if documentPath != "" {
		// Load the data
		log.Debug().Msgf("Loading data from: %v", documentPath)
		wg.Add(1)
		// Load the data into the collection
		go embdedData(ctx, documentPath, client, constants.TXT, &wg, errChan, collectionChan)

		select {
		case <-errChan:
			log.Error().Msgf("Error loading data: %v", err)
			return
		case <-collectionChan:
			log.Debug().Msg("Data loaded successfully")

		}

	}

	if queryString != "" {
		// queryString = "what is the difference between mirostat_tau and mirostat_eta?"
		log.Debug().Msgf("Querying with: %v", queryString)
		vectorQuery := []string{queryString}
		vectorChan := make(chan *chromago.QueryResults, 1)
		rankChan := make(chan string, 1)
		defer close(vectorChan)
		defer close(rankChan)

		collection, err := chromaclient.GetCollectionFromDb(ctx, client.Client, constants.HuggingFace, constants.HuggingFaceEmbedModel)
		if err != nil {
			log.Debug().Msgf("Error getting collection: %v\n", err)
			return
		}

		log.Debug().Msgf("Querying collection: %v", collection.Name)

		wg.Add(1)
		go queryVectorDB(ctx, collection, vectorQuery, &wg, errChan, vectorChan)

		wg.Add(1)
		go rerankQueryResults(ctx, vectorQuery, vectorChan, reRankClient, &wg, errChan, rankChan)

		// wait for all go routines to finish
		wg.Wait()

		// Get the final prompt and chat result
		var rankResult string
		select {
		case rankResult = <-rankChan:
		case <-ctx.Done():
			log.Error().Msgf("Context Timeout: %v", ctx.Err())
			return
		case <-errChan:
			log.Error().Msgf("Error querying vector db: %v", err)
			return
		}

		contentString := rankResult
		// Get the final prompt
		prompts, err := prompts.GetTemplate(constants.SystemPromptFile, queryString, contentString)
		if err != nil {
			log.Error().Msgf("Failed to get template: %v", err)

		}

		log.Debug().Msgf("Final Prompt: %v", prompts)
		provider := chat.NewChatManager(constants.OllamaChat, constants.OllamaChatModel, constants.OllamaUrl, "")
		provider.Chat(ctx, prompts)
	}

}

func embdedData(ctx context.Context, path string, client *chromaclient.ChromaClient, docType constants.DocType, wg *sync.WaitGroup, errChan chan<- error, collectionChan chan<- *chromago.Collection) {
	defer wg.Done()
	collection, err := documenthandler.VectorEmbedData(ctx, client,
		documenthandler.WithDocPath(path),
		documenthandler.WithDocType(docType))
	if err != nil {
		errChan <- err
	}
	collectionChan <- collection
}

func queryVectorDB(ctx context.Context, collection *chromago.Collection, query []string, wg *sync.WaitGroup, errChan chan<- error, vectorChan chan<- *chromago.QueryResults) {
	defer wg.Done()
	vectorResults, err := vectordbquery.QueryVectorDbWithOptions(ctx, collection, query)
	if err != nil {
		errChan <- err
		log.Error().Msgf("Failed to query vector db: %v", err)
	} else {
		vectorChan <- vectorResults
	}
}

func rerankQueryResults(ctx context.Context, query []string, vectorChan chan *chromago.QueryResults, reRankClient reranker.Reranker, wg *sync.WaitGroup, errChan chan<- error, rankChan chan<- string) {
	defer wg.Done()
	queryResults := <-vectorChan
	rerankResults, err := reRankClient.RerankQueryResult(ctx, query, queryResults.Documents[0])
	if err != nil {
		errChan <- err
		log.Error().Msgf("Failed to rerank query results: %v", err)
	} else {
		rankChan <- rerankResults
	}
}

func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		log.Error().Msgf("Error reading input: %v", err)
		return ""
	}
	return input[:len(input)-1] // Remove the newline character
}
