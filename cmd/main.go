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

// TODO : Refactor the code to split the main function into smaller functions
func main() {
	ctx := context.Background()
	// sometimes timeout happens while model is running on remote server switching to cancel context //TODO
	// ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Define the flags
	loadFlag := flag.Bool("load", false, "Load and embed the data in vectordb")
	queryFlag := flag.Bool("query", false, "Query the embedded data and rerank the results")
	grpcFlag := flag.Bool("grpc", false, "Query the embedded data and rerank the results using gRPC")

	// Parse the flags
	flag.Parse()

	// TODO : POC remove after incorporate this in the main code
	if *grpcFlag {
		query := []string{"what is Deep Learning?"}
		texts := []string{"Tomatos are fruits..", "Deep Learning is not...", "Deep learning is..."}
		grpcClient, _ := reranker.NewReRanker(ctx, constants.GRPC)
		result, err := grpcClient.RerankQueryResult(ctx, query, texts)
		if err != nil {
			log.Error().Msgf("Error initializing GrpcClient: %v\n", err)
			return
		}

		log.Debug().Msgf("Rerank Result: %v\n", result)

	}

	if !*loadFlag && !*queryFlag {
		fmt.Println("Please specify a flag: -load or -query")
		return
	}

	// Create a new reader to read from standard input
	reader := bufio.NewReader(os.Stdin)

	var documentPath string
	// If the load flag is set, load the data // TODO: prompt collection name
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
		// Load the data into the collection
		go embdedData(ctx, documentPath, client, constants.TXT, &wg, errChan, collectionChan)

		c := <-collectionChan
		log.Debug().Msgf("Loaded data to Collection: %v", c.Name)

	}

	if queryString != "" {
		// queryString = "what is the difference between mirostat_tau and mirostat_eta?"
		log.Debug().Msgf("Querying with: %v", queryString)
		// Query the collection with the query text
		vectorQuery := []string{queryString}
		vectorChan := make(chan *chromago.QueryResults, 1)
		rankChan := make(chan string, 1)
		defer close(vectorChan)
		defer close(rankChan)

		// c := <-collectionChan
		c, err := chromaclient.GetCollectionFromDb(ctx, client.Client, constants.HuggingFace, constants.HuggingFaceEmbedModel)
		if err != nil {
			log.Debug().Msgf("Error getting collection: %v\n", err)
			return
		}

		log.Debug().Msgf("Querying collection: %v", c.Name)

		wg.Add(1)
		// Query the vector db
		go queryVectorDB(ctx, c, vectorQuery, &wg, errChan, vectorChan)

		// Initialize the ReRank client
		reRankClient, err := reranker.NewReRanker(ctx, constants.GRPC)
		if err != nil {
			log.Error().Msgf("Error initializing ReRank Client: %v\n", err)
			return
		}
		wg.Add(1)
		// Rerank the query results
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
		}

		contentString := rankResult
		// Get the final prompt
		prompts, err := prompts.GetTemplate(constants.SystemPromptFile, queryString, contentString)
		if err != nil {
			log.Error().Msgf("Failed to get template: %v", err)

		}

		log.Debug().Msgf("Final Prompt: %v", prompts)

		// Chat with the user using the final prompt and content
		chat.ChatOllama(ctx, prompts)
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
