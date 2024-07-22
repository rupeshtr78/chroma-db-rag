package db

import (
	"chroma-db/internal/chromaclient"
	"chroma-db/internal/constants"
	"chroma-db/internal/documents"
	ollamamodel "chroma-db/internal/ollama"
	"chroma-db/internal/vectordb"
	"context"
	"log"

	"github.com/amikos-tech/chroma-go/types"
)

func RunVectorDb(ctx context.Context) {
	// Get the chroma client
	client, err := chromaclient.GetChromaClient(ctx, constants.ChromaUrl)
	if err != nil {
		log.Default().Println(err)
		return
	}

	// Get or create the tenant
	t, err := chromaclient.GetOrCreateTenant(ctx, client, constants.TenantName)
	if err != nil {
		log.Default().Println(err)
		return
	}

	// Get or create the database
	d, err := chromaclient.GetOrCreateDatabase(ctx, client, constants.Database, t.Name)
	if err != nil {
		log.Default().Println(err)
		return
	}
	log.Default().Printf("Database: %v\n", *d.Name)

	// Get the ollama embedding function
	ollamaEmbedFn, err := ollamamodel.GetOllamaEmbedding(constants.OllamaUrl, constants.OllamaModel)
	if err != nil {
		log.Default().Println(err)
		return
	}

	// Create a new store
	store, err := chromaclient.CreateChromaStore(ctx,
		constants.ChromaUrl,
		constants.Namespace,
		ollamaEmbedFn,
		types.DistanceFunction(constants.DistanceFn))
	if err != nil {
		log.Default().Println(err)
		return
	}

	// err4 := store.RemoveCollection()
	// if err4 != nil {
	// 	log.Default().Println(err4)
	// 	return
	// }

	// // get the documents from the pdf
	pdfDocs, err := documents.PdfToDocument(ctx, "text/Model Params.pdf")
	if err != nil {
		log.Default().Println(err)
		return
	}

	// Add the documents to the store
	s, err3 := store.AddDocuments(ctx, pdfDocs)
	if err3 != nil {
		log.Default().Println(err3)
		return
	}
	log.Default().Printf("Added %v documents\n", len(s))

	// Get the list of all the available collections
	collections, err2 := client.ListCollections(ctx)
	if err2 != nil {
		log.Default().Println(err2)
	}

	// Print the list of databases

	for _, col := range collections {
		log.Default().Printf("Collection: %v\n", col.Name)
		log.Default().Printf("Database: %v\n", col.Database)
		log.Default().Printf("Tenant: %v\n", col.Tenant)
	}

	// // Search the store
	queryString := "what is mirostat_tau?"
	docs, err := vectordb.SearchVectorDb(ctx,
		store,
		queryString,
		3,
		constants.Namespace)

	if err != nil {
		log.Default().Println(err)
		return
	}

	log.Default().Printf("Found %v documents\n", len(docs))
	for _, doc := range docs {
		log.Default().Printf("Document: %v\n", doc.PageContent)
		log.Default().Printf("Document: %v\n", doc.Metadata)
		log.Default().Printf("Document: %v\n", doc.Score)

	}
}
