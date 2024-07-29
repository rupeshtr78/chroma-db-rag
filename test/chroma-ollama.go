package test

import (
	"context"
	"fmt"
	"log"
	"strings"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

var CHROMA_URL = "http://0.0.0.0:8070"

var namespace = "chroma-ollama"

var score_threshold float32 = 0.6

func ChromaEmbedder() {
	// Create a new Ollama embedder.
	// ollama pull nomic-embed-text ,ollama serve mxbai-embed-large
	ollamaEmbeder := GetOllamaEmbedding("mxbai-embed-large")

	// clientOpts := vectorstores.Options{NameSpace: namespace}
	// create the client connection and confirm that we can access the server with it
	chromaClient, err := chromago.NewClient(CHROMA_URL)
	if err != nil {
		log.Fatalf("new client: %v\n", err)
	}

	if _, errHb := chromaClient.Heartbeat(context.Background()); errHb != nil {
		log.Fatalf("heartbeat: %v\n", errHb)
	}

	// get collection
	ollamaCollection, err := chromaClient.GetCollection(context.Background(), namespace, nil)
	if err != nil {
		log.Fatalf("get collection: %v\n", err)
	}

	fmt.Printf("Got collection: %v\n", ollamaCollection.Name)

	// delete collection
	c, err := chromaClient.DeleteCollection(context.Background(), namespace)
	if err != nil {
		log.Fatalf("error delete collection: %v\n", err)
	}
	fmt.Printf("deleted collection: %v\n", c.Name)

	// Create a new Chroma vector store.
	store := CreateChromaStore(ollamaEmbeder, namespace)

	// Add sample data to the vector store.
	AddSampleData(store, namespace)

	ctx := context.TODO()

	// Create example cases.
	exampleCases := SampleQuery()

	// run the example cases
	results := make([][]schema.Document, len(exampleCases))

	// query_results := make([]chromago.QueryResults, len(exampleCases))
	// count collection
	coll, err := GetCollection(ctx, chromaClient)
	if err != nil {
		fmt.Println("Error fecthing collection count")
	}

	count, err := coll.Count(ctx)
	if err != nil {
		log.Fatalf("count: %v\n", err)
	}

	fmt.Printf("Collecton count: %v\n", count)
	query := make([]string, len(exampleCases))

	for _, ec := range exampleCases {
		query = append(query, ec.query)

		qr, err := coll.Query(ctx,
			query,
			1,
			nil,
			nil,
			nil)
		if err != nil {
			log.Fatalf("query1: %v\n", err)
		}

		if len(qr.Documents) > 0 && len(qr.Documents[0]) > 0 {
			fmt.Printf("qr: %v\n", qr.Documents[0][0]) // this should result
		} else {
			log.Fatalf("No documents returned")
		}

	}

	// print out the results of the run

	for ecI, ec := range exampleCases {
		docs, errSs := store.SimilaritySearch(ctx, ec.query, ec.numDocuments, ec.options...)
		if errSs != nil {
			log.Fatalf("query1: %v\n", errSs)
		}
		results[ecI] = docs
	}

	// print out the results of the run
	fmt.Printf("Results:\n")
	for ecI, ec := range exampleCases {
		texts := make([]string, len(results[ecI]))
		for docI, doc := range results[ecI] {
			texts[docI] = doc.PageContent
		}
		fmt.Printf("%d. case: %s\n", ecI+1, ec.name)
		fmt.Printf("    result: %s\n", strings.Join(texts, ", "))
	}
}

type exampleCase struct {
	name         string
	query        string
	numDocuments int
	options      []vectorstores.Option
}

func SampleQuery() []exampleCase {

	type filter = map[string]any

	exampleCases := []exampleCase{
		{
			name:         "Up to 5 Cities in Japan",
			query:        "Which of these are cities are located in Japan?",
			numDocuments: 5,
			options: []vectorstores.Option{
				vectorstores.WithScoreThreshold(score_threshold),
			},
		},
		{
			name:         "A City in South America",
			query:        "Which of these are cities are located in South America?",
			numDocuments: 1,
			options: []vectorstores.Option{
				vectorstores.WithScoreThreshold(score_threshold),
			},
		},
		{
			name:         "Large Cities in South America",
			query:        "Which of these are cities are located in South America?",
			numDocuments: 100,
			options: []vectorstores.Option{
				vectorstores.WithFilters(filter{
					"$and": []filter{
						{"area": filter{"$gte": 1000}},
						{"population": filter{"$gte": 13}},
					},
				}),
				vectorstores.WithScoreThreshold(score_threshold),
			},
		},
	}
	return exampleCases
}

func GetCollection(ctx context.Context, chromaClient *chromago.Client) (*chromago.Collection, error) {
	ollamaCollection, err := chromaClient.GetCollection(context.Background(), namespace, nil)
	if err != nil {
		log.Fatalf("get collection: %v\n", err)
		return nil, err
	}

	return ollamaCollection, nil
}

func getChromeStore() (chroma.Store, error) {
	store, err := chroma.New(
		chroma.WithChromaURL(CHROMA_URL),
	)
	if err != nil {
		log.Fatalf("new: %v\n", err)
		return chroma.Store{}, err
	}
	return store, err
}

func CreateChromaStore(ollamaEmbeder *embeddings.EmbedderImpl, namespace string) chroma.Store {
	store, errNs := chroma.New(
		chroma.WithChromaURL(CHROMA_URL),
		chroma.WithEmbedder(ollamaEmbeder),
		chroma.WithDistanceFunction("cosine"), // l2, cosine, ip
		// chroma.WithNameSpace(uuid.New().String()),
		chroma.WithNameSpace(namespace),
	)
	if errNs != nil {
		log.Fatalf("new: %v\n", errNs)
	}
	return store
}

func GetOllamaEmbedding(model string) *embeddings.EmbedderImpl {
	ollamaLLM, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		log.Fatal(err)
	}
	ollamaEmbeder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		log.Fatal(err)
	}
	return ollamaEmbeder
}

func AddSampleData(store chroma.Store, namespace string) {

	// Add documents to the vector store.
	_, errAd := store.AddDocuments(context.Background(), docs,
		vectorstores.WithNameSpace(namespace),
	)
	if errAd != nil {
		log.Fatalf("AddDocument: %v\n", errAd)
	}
}

type meta = map[string]any

var docs = []schema.Document{
	{PageContent: "Tokyo", Metadata: meta{"population": 9.7, "area": 622}},
	{PageContent: "Kyoto", Metadata: meta{"population": 1.46, "area": 828}},
	{PageContent: "Hiroshima", Metadata: meta{"population": 1.2, "area": 905}},
	{PageContent: "Kazuno", Metadata: meta{"population": 0.04, "area": 707}},
	{PageContent: "Nagoya", Metadata: meta{"population": 2.3, "area": 326}},
	{PageContent: "Toyota", Metadata: meta{"population": 0.42, "area": 918}},
	{PageContent: "Fukuoka", Metadata: meta{"population": 1.59, "area": 341}},
	{PageContent: "Paris", Metadata: meta{"population": 11, "area": 105}},
	{PageContent: "London", Metadata: meta{"population": 9.5, "area": 1572}},
	{PageContent: "Santiago", Metadata: meta{"population": 6.9, "area": 641}},
	{PageContent: "Buenos Aires", Metadata: meta{"population": 15.5, "area": 203}},
	{PageContent: "Rio de Janeiro", Metadata: meta{"population": 13.7, "area": 1200}},
	{PageContent: "Sao Paulo", Metadata: meta{"population": 22.6, "area": 1523}},
}
